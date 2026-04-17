package system

import (
	"github.com/gin-gonic/gin"

	"github.com/huuhoait/gin-vue-admin/server/middleware"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	{
		// Login is the primary bruteforce / credential-stuffing target, so
		// apply a tighter per-IP limiter than the global one.
		baseRouter.POST("login", middleware.LoginLimit(), baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
	}
	return baseRouter
}
