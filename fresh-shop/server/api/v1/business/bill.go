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

// GetBillList 获取账单列表
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
	response.OkWithDetailed(gin.H{
		"list":  list,
		"total": total,
	}, "查询成功", c)
}

// GenerateBill 生成账单
// POST /api/business/bill/generate
// Body: { "year": 2026, "month": 9 }
func (api *BillApi) GenerateBill(c *gin.Context) {
	var req struct {
		Year  int `json:"year" binding:"required"`
		Month int `json:"month" binding:"required"`
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

	if err := billService.GenerateBill(companyId, req.Year, req.Month); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("账单生成成功", c)
}

// GetBillDetail 获取账单详情
// GET /api/business/bill/:id
func (api *BillApi) GetBillDetail(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	bill, err := billService.GetBillDetail(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(bill, "查询成功", c)
}

// UpdateBillStatus 更新账单状态
// PUT /api/business/bill/:id/status
func (api *BillApi) UpdateBillStatus(c *gin.Context) {
	var req struct {
		ID     uint `uri:"id" binding:"required"`
		Status int  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := billService.UpdateBillStatus(req.ID, req.Status); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("状态更新成功", c)
}
