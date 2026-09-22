package middleware

import (
	"fresh-shop/server/model/system/request"
	"github.com/gin-gonic/gin"
)

// TenantContext 租户上下文
type TenantContext struct {
	CompanyID uint
	UserID    uint
	UserType  uint // 1超管 2公司管理员 3商户
}

// TenantInject 租户上下文注入中间件
// 在 JWTAuth 之后调用，从 claims 中提取租户信息
func TenantInject() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.Next()
			return
		}

		sysClaims := claims.(*request.CustomClaims)
		ctx := &TenantContext{
			UserID:    sysClaims.ID,
			UserType:  sysClaims.UserType,
			CompanyID: sysClaims.CompanyID,
		}
		c.Set("tenant", ctx)
		c.Next()
	}
}

// GetTenant 获取租户上下文
func GetTenant(c *gin.Context) *TenantContext {
	if ctx, exists := c.Get("tenant"); exists {
		return ctx.(*TenantContext)
	}
	return nil
}

// GetCompanyID 获取当前公司ID
func GetCompanyID(c *gin.Context) uint {
	if ctx := GetTenant(c); ctx != nil {
		return ctx.CompanyID
	}
	return 0
}

// IsSuperAdmin 是否超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	if ctx := GetTenant(c); ctx != nil {
		return ctx.UserType == 1
	}
	return false
}

// IsCompanyAdmin 是否公司管理员
func IsCompanyAdmin(c *gin.Context) bool {
	if ctx := GetTenant(c); ctx != nil {
		return ctx.UserType == 2
	}
	return false
}
