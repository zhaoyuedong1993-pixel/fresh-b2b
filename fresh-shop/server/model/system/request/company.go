package request

// CompanyByIdReq 根据ID获取公司请求
type CompanyByIdReq struct {
	ID uint `json:"id" uri:"id" binding:"required"`
}

// CompanyCreateReq 创建公司请求
type CompanyCreateReq struct {
	Name        string `json:"name" binding:"required"`        // 公司名称
	Address     string `json:"address"`                       // 公司地址
	Phone       string `json:"phone" binding:"required"`      // 联系电话
	ContactName string `json:"contactName"`                   // 负责人姓名
	Username    string `json:"username" binding:"required"`   // 管理员账号
	Password    string `json:"password" binding:"required"`   // 管理员密码
}

// CompanyUpdateReq 更新公司请求
type CompanyUpdateReq struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name"`                         // 公司名称
	Address     string `json:"address"`                     // 公司地址
	Phone       string `json:"phone"`                       // 联系电话
	ContactName string `json:"contactName"`                 // 负责人姓名
	Status      int    `json:"status"`                      // 状态
}

// CompanyAuditReq 公司审核请求
type CompanyAuditReq struct {
	ID       uint   `json:"id" binding:"required"`          // 公司ID
	Pass     bool   `json:"pass"`                          // 是否通过
	Remark   string `json:"remark"`                         // 备注
	Password string `json:"password"`                      // 管理员密码（审核通过时设置）
}

// CompanyRegisterReq 小程序申请入驻平台请求
type CompanyRegisterReq struct {
	Name        string `json:"name" binding:"required"`     // 公司名称
	Address     string `json:"address"`                    // 公司地址
	ContactName string `json:"contactName"`                // 负责人姓名
	Phone       string `json:"phone" binding:"required"`  // 手机号
	Password    string `json:"password" binding:"required"` // 密码
}

// CompanyJoinReq 小程序申请加入公司请求
type CompanyJoinReq struct {
	InvitationCode string `json:"invitationCode" binding:"required"` // 邀请码
	Phone          string `json:"phone" binding:"required"`           // 手机号
	Password       string `json:"password" binding:"required"`       // 密码
}
