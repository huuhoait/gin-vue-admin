package middleware

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// CasbinHandler interceptor
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		// get request PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// get request method
		act := c.Request.Method
		// get user's role
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin() // check if policy exists
		if e == nil {
			response.FailWithDetailed(gin.H{}, "Authorization subsystem unavailable", c)
			c.Abort()
			return
		}
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "Insufficient permissions", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
