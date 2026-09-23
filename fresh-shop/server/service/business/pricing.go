package business

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"fresh-shop/server/global"
)

// ParsedPrice 解析出的价格数据
type ParsedPrice struct {
	Name    string  `json:"name"`    // 商品名称
	Price   float64 `json:"price"`   // 价格
	Unit    string  `json:"unit"`    // 单位
	Updated bool    `json:"updated"` // 是否更新成功
}

// AIParseRequest AI解析请求
type AIParseRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// AIParseResponse AI解析响应
type AIParseResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// PricingService AI调价服务
type PricingService struct{}

// NewPricingService 创建定价服务
func NewPricingService() *PricingService {
	return &PricingService{}
}

// CalculateSellPrice 计算商户价
func (s *PricingService) CalculateSellPrice(costPrice, markupRate float64) float64 {
	return costPrice * (1 + markupRate/100)
}

// BatchUpdateCostPrice 批量更新采购价
func (s *PricingService) BatchUpdateCostPrice(companyId uint, updates map[uint]float64) error {
	for goodsId, costPrice := range updates {
		global.DB.Table("shop_goods").
			Where("id = ? AND company_id = ?", goodsId, companyId).
			Updates(map[string]interface{}{"cost_price": costPrice})
	}
	return nil
}

// ApplyPricing 批量应用价格（根据采购价+加价比例计算商户价）
func (s *PricingService) ApplyPricing(companyId uint) error {
	var company struct {
		MarkupRate float64
	}
	global.DB.Table("sys_company").Select("markup_rate").Where("id = ?", companyId).First(&company)

	var goods []struct {
		ID        uint
		CostPrice float64
	}
	global.DB.Table("shop_goods").
		Select("id, cost_price").
		Where("company_id = ? AND cost_price > 0", companyId).
		Find(&goods)

	for _, g := range goods {
		sellPrice := s.CalculateSellPrice(g.CostPrice, company.MarkupRate)
		global.DB.Table("shop_goods").Where("id = ?", g.ID).Update("price", sellPrice)
	}
	return nil
}

// GetMarkupRate 获取加价比例
func (s *PricingService) GetMarkupRate(companyId uint) (float64, error) {
	var company struct {
		MarkupRate float64
	}
	err := global.DB.Table("sys_company").Select("markup_rate").Where("id = ?", companyId).First(&company).Error
	if err != nil {
		return 10, nil // 默认10%
	}
	return company.MarkupRate, nil
}

// BatchParseAndUpdate 解析文案并批量更新采购价
func (s *PricingService) BatchParseAndUpdate(companyId uint, rawText string) ([]ParsedPrice, error) {
	// 调用AI解析
	parsedItems, err := s.parseWithAI(rawText)
	if err != nil {
		return nil, fmt.Errorf("AI解析失败: %v", err)
	}

	// 批量更新采购价
	for i, item := range parsedItems {
		// 查找匹配的商品
		goodsID := s.findGoodsByName(companyId, item.Name)
		if goodsID > 0 {
			global.DB.Table("shop_goods").
				Where("id = ? AND company_id = ?", goodsID, companyId).
				Update("cost_price", item.Price)
			parsedItems[i].Updated = true
		}
	}

	return parsedItems, nil
}

// parseWithAI 使用AI解析文案
func (s *PricingService) parseWithAI(text string) ([]ParsedPrice, error) {
	if !global.Config.AI.Enabled || global.Config.AI.APIKey == "" {
		// 如果AI未启用，使用正则解析
		return s.parseWithRegex(text)
	}

	prompt := fmt.Sprintf(`你是一个生鲜商品价格解析助手。请从以下文本中提取商品名称和价格，返回JSON数组格式。

要求：
1. 商品名称只提取核心词，如"莲藕"、"豆芽"、"土豆"，去掉"今日"等时间词
2. 价格提取数字，如"1.5"对应1.5元/斤
3. 单位统一为"斤"
4. 如果价格带"块"字，转换为元
5. 如果没有明确单位，默认是斤

示例输入："今日莲藕1.5块一斤，豆芽1.5块一斤"
示例输出：[{"name":"莲藕","price":1.5,"unit":"斤"},{"name":"豆芽","price":1.5,"unit":"斤"}]

输入文本：
%s

请直接返回JSON数组，不要有其他内容：`, text)

	// 调用AI API
	response, err := s.callAIAPI(prompt)
	if err != nil {
		// AI调用失败，降级到正则解析
		return s.parseWithRegex(text)
	}

	// 解析AI返回的JSON
	var items []ParsedPrice
	if err := json.Unmarshal([]byte(response), &items); err != nil {
		// 解析失败，降级到正则解析
		return s.parseWithRegex(text)
	}

	return items, nil
}

// callAIAPI 调用AI API
func (s *PricingService) callAIAPI(prompt string) (string, error) {
	reqBody := AIParseRequest{
		Model: global.Config.AI.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", global.Config.AI.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", global.Config.AI.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API调用失败: %d - %s", resp.StatusCode, string(body))
	}

	var aiResp AIParseResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", err
	}

	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("AI返回为空")
	}

	return aiResp.Choices[0].Message.Content, nil
}

// parseWithRegex 使用正则解析文案（降级方案）
func (s *PricingService) parseWithRegex(text string) ([]ParsedPrice, error) {
	var items []ParsedPrice

	// 清理文本
	text = strings.ReplaceAll(text, "今日", "")
	text = strings.ReplaceAll(text, "，", ",")
	text = strings.ReplaceAll(text, "。", ",")
	text = strings.ReplaceAll(text, "；", ",")
	text = strings.ReplaceAll(text, ",", ",")

	// 匹配模式：(商品)(价格)(单位)
	patterns := []string{
		`([\x{4e00}-\x{9fa5}]+)\s*([0-9]+\.?[0-9]*)\s*(?:块|元)(?:/|\s*)(?:斤|公斤|kg|KG)?`, // 1.5块一斤
		`([\x{4e00}-\x{9fa5}]+)\s*([0-9]+\.?[0-9]*)\s*(?:块|元)`,                                  // 1.5块
		`([\x{4e00}-\x{9fa5}]+)\s*([0-9]+\.?[0-9]*)\s*(?:斤|公斤|kg)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				name := strings.TrimSpace(match[1])
				priceStr := strings.TrimSpace(match[2])

				// 跳过明显不是商品的词
				if len(name) < 2 {
					continue
				}

				var price float64
				fmt.Sscanf(priceStr, "%f", &price)

				items = append(items, ParsedPrice{
					Name:    name,
					Price:   price,
					Unit:    "斤",
					Updated: false,
				})
			}
		}
	}

	return items, nil
}

// findGoodsByName 根据名称查找商品ID
func (s *PricingService) findGoodsByName(companyId uint, name string) uint {
	var goods struct {
		ID   uint
		Name string
	}

	// 模糊匹配商品名称
	searchPattern := "%" + name + "%"
	global.DB.Table("shop_goods").
		Select("id, name").
		Where("company_id = ? AND name LIKE ?", companyId, searchPattern).
		Where("status = 1").
		First(&goods)

	return goods.ID
}
