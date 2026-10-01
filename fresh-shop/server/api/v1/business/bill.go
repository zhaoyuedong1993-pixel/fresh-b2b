package business

import (
	"fresh-shop/server/middleware"
	"fresh-shop/server/model/common/response"
	"fresh-shop/server/service/business"
	"github.com/gin-gonic/gin"
)

// BillApi 账单API
type BillApi struct{}

var billService = &business.BillService{}

// GetBillList 获取账单列表（按月聚合）
// GET /api/business/bill/list?page=1&pageSize=10
func (api *BillApi) GetBillList(c *gin.Context) {
	var req struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	companyId := middleware.GetCompanyID(c)
	if companyId == 0 {
		response.FailWithMessage("无法获取公司信息", c)
		return
	}

	list, total := billService.GetBillList(companyId, req.Page, req.PageSize)
	response.OkWithData(gin.H{
		"list":  list,
		"total": total,
	}, c)
}

// GetBillDetail 获取账单详情
// GET /api/business/bill/detail?period=2026-09
func (api *BillApi) GetBillDetail(c *gin.Context) {
	period := c.Query("period")
	if period == "" {
		response.FailWithMessage("请提供账期", c)
		return
	}

	companyId := middleware.GetCompanyID(c)
	if companyId == 0 {
		response.FailWithMessage("无法获取公司信息", c)
		return
	}

	bill, err := billService.GetBillDetailByPeriod(companyId, period)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(bill, c)
}

// UpdateBillStatus 更新账单状态
// PUT /api/business/bill/status
func (api *BillApi) UpdateBillStatus(c *gin.Context) {
	var req struct {
		Period string `json:"period" binding:"required"`
		Status int    `json:"status" binding:"required"`
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

	if err := billService.UpdateBillStatusByPeriod(companyId, req.Period, req.Status); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("状态更新成功", c)
}
