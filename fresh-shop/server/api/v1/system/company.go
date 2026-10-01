package system

import (
	"fresh-shop/server/global"
	"fresh-shop/server/model/common/response"
	companyReq "fresh-shop/server/model/system/request"
	"fresh-shop/server/service"
	"fresh-shop/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type CompanyApi struct{}

var companyService = service.ServiceGroupApp.SystemServiceGroup.CompanyService

// CreateCompany 创建公司（超管）
func (companyApi *CompanyApi) CreateCompany(c *gin.Context) {
	var req companyReq.CompanyCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	company, err := companyService.CreateCompany(req)
	if err != nil {
		global.Log.Error("创建公司失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(company, "创建成功", c)
}

// GetCompanyList 获取公司列表
func (companyApi *CompanyApi) GetCompanyList(c *gin.Context) {
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			pageSize = v
		}
	}

	list, total, err := companyService.GetCompanyList(page, pageSize)
	if err != nil {
		global.Log.Error("获取公司列表失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GetCompany 获取单个公司
func (companyApi *CompanyApi) GetCompany(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("请选择公司", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage("无效的公司ID", c)
		return
	}

	company, err := companyService.GetCompany(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(company, "获取成功", c)
}

// UpdateCompany 更新公司
func (companyApi *CompanyApi) UpdateCompany(c *gin.Context) {
	var req companyReq.CompanyUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := companyService.UpdateCompany(req); err != nil {
		global.Log.Error("更新公司失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteCompany 删除公司
func (companyApi *CompanyApi) DeleteCompany(c *gin.Context) {
	var req companyReq.CompanyByIdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := companyService.DeleteCompany(req.ID); err != nil {
		global.Log.Error("删除公司失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetPendingCompanies 获取待审核入驻申请（公司审批）
func (companyApi *CompanyApi) GetPendingCompanies(c *gin.Context) {
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			pageSize = v
		}
	}

	list, total, err := companyService.GetPendingCompanies(page, pageSize)
	if err != nil {
		global.Log.Error("获取待审核列表失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// AuditCompany 审核入驻申请
func (companyApi *CompanyApi) AuditCompany(c *gin.Context) {
	var req companyReq.CompanyAuditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := companyService.AuditCompany(req); err != nil {
		global.Log.Error("审核失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("审核完成", c)
}

// RegisterCompany 小程序申请入驻平台
func (companyApi *CompanyApi) RegisterCompany(c *gin.Context) {
	var req companyReq.CompanyRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := companyService.RegisterCompany(req); err != nil {
		global.Log.Error("申请入驻失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("申请已提交，请等待审核", c)
}

// JoinCompany 小程序申请加入公司
func (companyApi *CompanyApi) JoinCompany(c *gin.Context) {
	var req companyReq.CompanyJoinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	if err := companyService.JoinCompany(req, userId); err != nil {
		global.Log.Error("申请加入公司失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("申请已提交，请等待管理员审核", c)
}

// GetCompanyByInvitationCode 通过邀请码获取公司信息
func (companyApi *CompanyApi) GetCompanyByInvitationCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("邀请码不能为空", c)
		return
	}

	company, err := companyService.GetCompanyByInvitationCode(code)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{
		"name":  company.Name,
		"phone": company.Phone,
	}, "获取成功", c)
}
