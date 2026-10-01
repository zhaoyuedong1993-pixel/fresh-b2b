package system

import (
	"fresh-shop/server/api/v1/system"
	"github.com/gin-gonic/gin"
)

type CompanyRouter struct{}

func (s *CompanyRouter) InitCompanyRouter(Router *gin.RouterGroup) {
	companyApi := new(system.CompanyApi)
	RouterGroup := Router.Group("company")
	{
		RouterGroup.POST("create", companyApi.CreateCompany)       // 创建公司
		RouterGroup.GET("list", companyApi.GetCompanyList)         // 公司列表
		RouterGroup.GET("detail", companyApi.GetCompany)           // 公司详情
		RouterGroup.POST("update", companyApi.UpdateCompany)        // 更新公司
		RouterGroup.POST("delete", companyApi.DeleteCompany)       // 删除公司
		RouterGroup.GET("pendingList", companyApi.GetPendingCompanies) // 待审核列表
		RouterGroup.POST("audit", companyApi.AuditCompany)         // 审核申请
		RouterGroup.GET("invitationCode", companyApi.GetCompanyByInvitationCode) // 通过邀请码获取公司
	}
}
