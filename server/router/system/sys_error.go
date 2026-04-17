package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysErrorRouter struct{}

// InitSysErrorRouter initialize Error Log route information
func (s *SysErrorRouter) InitSysErrorRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
    sysErrorRouter := Router.Group("sysError").Use(middleware.OperationRecord())
    sysErrorRouterWithoutRecord := Router.Group("sysError")
    sysErrorRouterWithoutAuth := PublicRouter.Group("sysError")
    {
        sysErrorRouter.DELETE("deleteSysError", sysErrorApi.DeleteSysError)           // deleteError Log
        sysErrorRouter.DELETE("deleteSysErrorByIds", sysErrorApi.DeleteSysErrorByIds) // batch delete errorsLog
        sysErrorRouter.PUT("updateSysError", sysErrorApi.UpdateSysError)              // updateError Log
        sysErrorRouter.GET("getSysErrorSolution", sysErrorApi.GetSysErrorSolution)    // TriggerError LogHandle
    }
    {
        sysErrorRouterWithoutRecord.GET("findSysError", sysErrorApi.FindSysError)       // get by IDError Log
        sysErrorRouterWithoutRecord.GET("getSysErrorList", sysErrorApi.GetSysErrorList) // getError LogList
    }
    {
        sysErrorRouterWithoutAuth.POST("createSysError", sysErrorApi.CreateSysError) // createError Log
    }
}
