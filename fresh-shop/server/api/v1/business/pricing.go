package business

import (
	"fresh-shop/server/middleware"
	"fresh-shop/server/model/common/response"
	"fresh-shop/server/service/business"
	"github.com/gin-gonic/gin"
)

// PricingApi AI调价API
type PricingApi struct{}

var pricingService = &business.PricingService{}

// BatchUpdateCostPrice 批量更新采购价
// POST /api/business/pricing/batch-cost-price
// Body: { "updates": { "goodsId": costPrice } }
func (api *PricingApi) BatchUpdateCostPrice(c *gin.Context) {
	var req struct {
		Updates map[uint]float64 `json:"updates"` // goodsId -> costPrice
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	companyId := middleware.GetCompanyID(c)
	if companyId == 0 {
		response.FailWithMessage("无法获取公司信息", c)
		return
	}

	if err := pricingService.BatchUpdateCostPrice(companyId, req.Updates); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("采购价更新成功", c)
}

// ApplyPricing 应用价格
// POST /api/business/pricing/apply
// 根据采购价和加价比例批量计算并更新商户价
func (api *PricingApi) ApplyPricing(c *gin.Context) {
	companyId := middleware.GetCompanyID(c)
	if companyId == 0 {
		response.FailWithMessage("无法获取公司信息", c)
		return
	}

	if err := pricingService.ApplyPricing(companyId); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("价格已生效", c)
}

// GetMarkupRate 获取当前公司的加价比例
// GET /api/business/pricing/markup-rate
func (api *PricingApi) GetMarkupRate(c *gin.Context) {
	companyId := middleware.GetCompanyID(c)
	if companyId == 0 {
		response.FailWithMessage("无法获取公司信息", c)
		return
	}

	rate, err := pricingService.GetMarkupRate(companyId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"markupRate": rate}, "查询成功", c)
}
