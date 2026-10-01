package system

import (
	"fresh-shop/server/global"
	"time"
)

// Company 公司/租户结构
type Company struct {
	global.DbModel
	Name           string  `json:"name" gorm:"column:name;not null;comment:公司名称;size:100;"`
	Contact        string  `json:"contact" gorm:"column:contact;comment:联系人;size:50;"`
	Phone          string  `json:"phone" gorm:"column:phone;comment:联系电话;size:20;"`
	Address        string  `json:"address" gorm:"column:address;comment:地址;size:200;"`
	CompanyType    string  `json:"companyType" gorm:"column:company_type;default:'monthly';comment:客户类型 monthly=月度 retail=零售;size:20;"`
	MarkupRate     float64 `json:"markupRate" gorm:"column:markup_rate;default:10.00;comment:加价比例(%);"`
	CutoffTime     string  `json:"cutoffTime" gorm:"column:cutoff_time;default:22:00;comment:截单时间;size:10;"`
	Status         int     `json:"status" gorm:"column:status;default:1;comment:状态 0禁用 1启用;"`

	// 新增字段（如果没有需要ALTER TABLE添加）
	InvitationCode string    `json:"invitationCode" gorm:"column:invitation_code;uniqueIndex;comment:邀请码;size:20;"`
	AdminUserId    uint      `json:"adminUserId" gorm:"column:admin_user_id;comment:管理员用户ID;"`
	ApplyTime      time.Time `json:"applyTime" gorm:"column:apply_time;comment:申请时间;"`
	AuditStatus    int8      `json:"auditStatus" gorm:"column:audit_status;default:0;comment:审核状态 0正常 1待审核入驻;"`
	AuditRemark    string    `json:"auditRemark" gorm:"column:audit_remark;comment:审核备注;size:255;"`
}

func (Company) TableName() string {
	return "sys_company"
}
