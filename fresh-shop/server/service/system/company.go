package system

import (
	"errors"
	"fresh-shop/server/global"
	"fresh-shop/server/model/system"
	systemReq "fresh-shop/server/model/system/request"
	"fresh-shop/server/utils"
	"time"
)

type CompanyService struct{}

// GenerateInvitationCode 生成邀请码
func generateInvitationCode() string {
	return utils.GenerateInviteCode(8)
}

// CreateCompany 创建公司（超管操作）
func (companyService *CompanyService) CreateCompany(req systemReq.CompanyCreateReq) (*system.Company, error) {
	// 检查公司名是否已存在
	var existing system.Company
	if err := global.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return nil, errors.New("公司名称已存在")
	}

	// 生成邀请码
	invitationCode := generateInvitationCode()

	// 创建公司记录
	company := system.Company{
		Name:           req.Name,
		Address:        req.Address,
		Phone:          req.Phone,
		Contact:        req.ContactName,
		InvitationCode: invitationCode,
		Status:         1,
		AuditStatus:    0, // 直接创建是已审核状态
		ApplyTime:      time.Now(),
	}

	if err := global.DB.Create(&company).Error; err != nil {
		return nil, errors.New("创建公司失败")
	}

	// 创建管理员用户
	adminUser := system.SysUser{
		Username:    req.Username,
		Password:    utils.BcryptHash(req.Password),
		NickName:    req.ContactName,
		Phone:       req.Phone,
		Enable:      1,
		LoginTime:   time.Now(),
		ApplyTime:   time.Now(),
		AuditStatus: 1, // 已审核
		CompanyId:   company.ID,
	}

	if err := global.DB.Create(&adminUser).Error; err != nil {
		// 回滚公司记录
		global.DB.Delete(&company)
		return nil, errors.New("创建管理员账号失败")
	}

	// 更新公司的管理员ID
	company.AdminUserId = adminUser.ID
	global.DB.Save(&company)

	// 设置管理员角色为公司管理员
	authLink := system.SysUserAuthority{
		SysUserId:                 adminUser.ID,
		SysAuthorityAuthorityId: 999, // 公司管理员角色
	}
	global.DB.Create(&authLink)

	return &company, nil
}

// GetCompanyList 获取公司列表
func (companyService *CompanyService) GetCompanyList(page, pageSize int) ([]system.Company, int64, error) {
	var list []system.Company
	var total int64

	db := global.DB.Model(&system.Company{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetCompany 获取单个公司
func (companyService *CompanyService) GetCompany(id uint) (*system.Company, error) {
	var company system.Company
	if err := global.DB.First(&company, id).Error; err != nil {
		return nil, errors.New("公司不存在")
	}
	return &company, nil
}

// UpdateCompany 更新公司
func (companyService *CompanyService) UpdateCompany(req systemReq.CompanyUpdateReq) error {
	var company system.Company
	if err := global.DB.First(&company, req.ID).Error; err != nil {
		return errors.New("公司不存在")
	}

	updates := map[string]interface{}{
		"name":          req.Name,
		"address":       req.Address,
		"phone":         req.Phone,
		"contact_name": req.ContactName,
		"status":        req.Status,
	}

	return global.DB.Model(&company).Updates(updates).Error
}

// DeleteCompany 删除公司
func (companyService *CompanyService) DeleteCompany(id uint) error {
	var company system.Company
	if err := global.DB.First(&company, id).Error; err != nil {
		return errors.New("公司不存在")
	}

	// 删除公司下的所有用户
	global.DB.Where("company_id = ?", id).Delete(&system.SysUser{})

	return global.DB.Delete(&company).Error
}

// GetPendingCompanies 获取待审核入驻申请（超管审批）
func (companyService *CompanyService) GetPendingCompanies(page, pageSize int) ([]system.Company, int64, error) {
	var list []system.Company
	var total int64

	db := global.DB.Model(&system.Company{}).Where("audit_status = 1")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// AuditCompany 审核入驻申请
func (companyService *CompanyService) AuditCompany(req systemReq.CompanyAuditReq) error {
	var company system.Company
	if err := global.DB.First(&company, req.ID).Error; err != nil {
		return errors.New("申请不存在")
	}

	if company.AuditStatus != 1 {
		return errors.New("该申请不是待审核状态")
	}

	if req.Pass {
		// 生成邀请码
		invitationCode := generateInvitationCode()

		// 创建管理员用户
		adminUser := system.SysUser{
			Username:    company.Phone,
			Password:   utils.BcryptHash(req.Password),
			NickName:   company.Contact,
			Phone:      company.Phone,
			Enable:     1,
			LoginTime:  time.Now(),
			ApplyTime:  time.Now(),
			AuditStatus: 1,
			CompanyId:  company.ID,
		}

		if err := global.DB.Create(&adminUser).Error; err != nil {
			return errors.New("创建管理员账号失败")
		}

		// 设置管理员角色
		authLink := system.SysUserAuthority{
			SysUserId:                 adminUser.ID,
			SysAuthorityAuthorityId: 999,
		}
		global.DB.Create(&authLink)

		// 更新公司状态
		return global.DB.Model(&company).Updates(map[string]interface{}{
			"audit_status":    0,
			"invitation_code": invitationCode,
			"admin_user_id":   adminUser.ID,
		}).Error
	} else {
		// 拒绝
		return global.DB.Model(&company).Updates(map[string]interface{}{
			"audit_status": 2,
			"audit_remark": req.Remark,
		}).Error
	}
}

// RegisterCompany 小程序申请入驻平台
func (companyService *CompanyService) RegisterCompany(req systemReq.CompanyRegisterReq) error {
	// 创建待审核公司记录
	company := system.Company{
		Name:        req.Name,
		Address:     req.Address,
		Phone:       req.Phone,
		Contact:     req.ContactName,
		Status:      1,
		AuditStatus: 1, // 待审核
		ApplyTime:   time.Now(),
	}

	return global.DB.Create(&company).Error
}

// JoinCompany 小程序申请加入公司
func (companyService *CompanyService) JoinCompany(req systemReq.CompanyJoinReq, userId uint) error {
	// 查找公司
	var company system.Company
	if err := global.DB.Where("invitation_code = ?", req.InvitationCode).First(&company).Error; err != nil {
		return errors.New("邀请码无效")
	}

	// 更新用户状态为待审核
	return global.DB.Model(&system.SysUser{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"audit_status": 0,
		"company_id":    company.ID,
	}).Error
}

// GetCompanyByInvitationCode 通过邀请码获取公司
func (companyService *CompanyService) GetCompanyByInvitationCode(code string) (*system.Company, error) {
	var company system.Company
	if err := global.DB.Where("invitation_code = ?", code).First(&company).Error; err != nil {
		return nil, errors.New("邀请码无效")
	}
	return &company, nil
}
