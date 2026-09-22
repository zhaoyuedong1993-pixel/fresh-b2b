package business

import (
	"fresh-shop/server/api/v1/business"
	"fresh-shop/server/middleware"
	"github.com/gin-gonic/gin"
)

// BusinessRouter 业务路由组
type BusinessRouter struct{}

// InitBusinessRouter 初始化业务路由
func (s *BusinessRouter) InitBusinessRouter(Router *gin.RouterGroup) {
	// 定价路由
	pricingRouter := Router.Group("pricing").Use(middleware.OperationRecord())
	{
		pricingRouter.POST("batch-cost-price", business.PricingApi{}.BatchUpdateCostPrice)
		pricingRouter.POST("apply", business.PricingApi{}.ApplyPricing)
		pricingRouter.GET("markup-rate", business.PricingApi{}.GetMarkupRate)
	}

	// 账单路由
	billRouter := Router.Group("bill").Use(middleware.OperationRecord())
	{
		billRouter.GET("list", business.BillApi{}.GetBillList)
		billRouter.POST("generate", business.BillApi{}.GenerateBill)
		billRouter.GET("/:id", business.BillApi{}.GetBillDetail)
		billRouter.PUT("/:id/status", business.BillApi{}.UpdateBillStatus)
	}
}
