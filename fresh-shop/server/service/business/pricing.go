package business

import (
	"fresh-shop/server/global"
)

// PricingService AI调价服务
type PricingService struct{}

// CalculateSellPrice 计算商户价
// 公式: 商户价 = 采购价 × (1 + 加价比例/100)
func (s *PricingService) CalculateSellPrice(costPrice, markupRate float64) float64 {
	return costPrice * (1 + markupRate/100)
}

// BatchUpdateCostPrice 批量更新采购价
// updates: map[goodsId]costPrice
func (s *PricingService) BatchUpdateCostPrice(companyId uint, updates map[uint]float64) error {
	for goodsId, costPrice := range updates {
		if costPrice <= 0 {
			continue
		}
		global.GlobalDb.Table("shop_goods").
			Where("id = ? AND company_id = ?", goodsId, companyId).
			Update("cost_price", costPrice)
	}
	return nil
}

// ApplyPricing 批量应用价格
// 根据采购价和公司加价比例，计算并更新商户价
func (s *PricingService) ApplyPricing(companyId uint) error {
	// 获取公司的加价比例
	var company struct {
		MarkupRate float64
	}
	if err := global.GlobalDb.Table("sys_company").
		Select("markup_rate").
		Where("id = ?", companyId).
		First(&company).Error; err != nil {
		return err
	}

	// 获取所有有采购价的商品
	var goodsList []struct {
		ID        uint
		CostPrice float64
	}
	if err := global.GlobalDb.Table("shop_goods").
		Select("id, cost_price").
		Where("company_id = ? AND cost_price > 0", companyId).
		Find(&goodsList).Error; err != nil {
		return err
	}

	// 批量更新商户价
	for _, g := range goodsList {
		sellPrice := s.CalculateSellPrice(g.CostPrice, company.MarkupRate)
		global.GlobalDb.Table("shop_goods").
			Where("id = ?", g.ID).
			Update("price", sellPrice)
	}

	return nil
}

// GetMarkupRate 获取公司的加价比例
func (s *PricingService) GetMarkupRate(companyId uint) (float64, error) {
	var company struct {
		MarkupRate float64
	}
	if err := global.GlobalDb.Table("sys_company").
		Select("markup_rate").
		Where("id = ?", companyId).
		First(&company).Error; err != nil {
		return 10.0, err // 默认10%
	}
	return company.MarkupRate, nil
}
