package middleware

import (
	"strconv"
	"strings"

	"fresh-shop/server/global"
	"fresh-shop/server/model/common/response"
	"fresh-shop/server/service"
	"fresh-shop/server/utils"
	"github.com/gin-gonic/gin"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		// 获取请求路径和方法
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.Config.System.RouterPrefix)
		act := c.Request.Method
		sub := strconv.Itoa(int(waitUse.AuthorityId))

		// 超级管理员(888)跳过权限检查
		if waitUse.AuthorityId == 888 {
			c.Next()
			return
		}

		if global.Config.System.Env != "develop" {
			e := casbinService.Casbin() // 判断策略中是否存在
			success, _ := e.Enforce(sub, obj, act)
			if !success {
				response.FailWithDetailed(gin.H{}, "权限不足", c)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
