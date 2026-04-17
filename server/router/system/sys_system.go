package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysRouter struct{}

func (s *SysRouter) InitSystemRouter(Router *gin.RouterGroup) {
	sysRouter := Router.Group("system").Use(middleware.OperationRecord())
	sysRouterWithoutRecord := Router.Group("system")

	{
		sysRouter.POST("setSystemConfig", systemApi.SetSystemConfig) // setconfigurationFilecontent
		sysRouter.POST("reloadSystem", systemApi.ReloadSystem)       // RepeatstartService
	}
	{
		sysRouterWithoutRecord.POST("getSystemConfig", systemApi.GetSystemConfig) // getconfigurationFilecontent
		sysRouterWithoutRecord.POST("getServerInfo", systemApi.GetServerInfo)     // get server info
	}
}
