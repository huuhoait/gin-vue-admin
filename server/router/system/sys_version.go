package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysVersionRouter struct{}

// InitSysVersionRouter initialize Version Management route information
func (s *SysVersionRouter) InitSysVersionRouter(Router *gin.RouterGroup) {
	sysVersionRouter := Router.Group("sysVersion").Use(middleware.OperationRecord())
	sysVersionRouterWithoutRecord := Router.Group("sysVersion")
	{
		sysVersionRouter.DELETE("deleteSysVersion", sysVersionApi.DeleteSysVersion)           // deleteVersion Management
		sysVersionRouter.DELETE("deleteSysVersionByIds", sysVersionApi.DeleteSysVersionByIds) // batch delete versionsmanagement
		sysVersionRouter.POST("exportVersion", sysVersionApi.ExportVersion)                   // export versionData
		sysVersionRouter.POST("importVersion", sysVersionApi.ImportVersion)                   // import versionData
	}
	{
		sysVersionRouterWithoutRecord.GET("findSysVersion", sysVersionApi.FindSysVersion)           // get by IDVersion Management
		sysVersionRouterWithoutRecord.GET("getSysVersionList", sysVersionApi.GetSysVersionList)     // getVersion ManagementList
		sysVersionRouterWithoutRecord.GET("downloadVersionJson", sysVersionApi.DownloadVersionJson) // download versionJSONData
	}
}
