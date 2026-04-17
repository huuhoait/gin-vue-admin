package system

import (
	"github.com/gin-gonic/gin"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
)

type GDPRRouter struct{}

func (r *GDPRRouter) InitGDPRRouter(Router *gin.RouterGroup) {
	// Self-service routes: authenticated user acts on their own data
	gdprSelf := Router.Group("user/gdpr").Use(middleware.OperationRecord())
	{
		gdprSelf.GET("export", gdprApiVar.ExportMyData)
		gdprSelf.DELETE("erase", gdprApiVar.EraseMyData)
	}
	// Admin route: admin acts on another user's data
	gdprAdmin := Router.Group("user/gdpr").Use(middleware.OperationRecord())
	{
		gdprAdmin.DELETE("adminErase", gdprApiVar.AdminEraseUser)
	}
}
