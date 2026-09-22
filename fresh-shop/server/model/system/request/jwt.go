package request

import (
	jwt "github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
	OpenId      string
	AuditStatus int8
	CompanyID   uint `json:"company_id"` // 公司ID
	UserType    uint `json:"user_type"` // 用户类型: 1超管 2公司管理员 3商户
}
