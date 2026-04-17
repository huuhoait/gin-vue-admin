package system

import (
	"github.com/huuhoait/gin-vue-admin/server/api/v1"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type LoginLogRouter struct{}

func (s *LoginLogRouter) InitLoginLogRouter(Router *gin.RouterGroup) {
	loginLogRouter := Router.Group("sysLoginLog").Use(middleware.OperationRecord())
	loginLogRouterWithoutRecord := Router.Group("sysLoginLog")
	sysLoginLogApi := v1.ApiGroupApp.SystemApiGroup.LoginLogApi
	{
		loginLogRouter.DELETE("deleteLoginLog", sysLoginLogApi.DeleteLoginLog)           // delete login log
		loginLogRouter.DELETE("deleteLoginLogByIds", sysLoginLogApi.DeleteLoginLogByIds) // batch delete login logs
	}
	{
		loginLogRouterWithoutRecord.GET("findLoginLog", sysLoginLogApi.FindLoginLog)       // get by IDLogin Log(Details)
		loginLogRouterWithoutRecord.GET("getLoginLogList", sysLoginLogApi.GetLoginLogList) // get login log list
	}
}
