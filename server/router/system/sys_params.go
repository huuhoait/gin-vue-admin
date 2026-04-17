package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysParamsRouter struct{}

// InitSysParamsRouter initialize Parameter route information
func (s *SysParamsRouter) InitSysParamsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sysParamsRouter := Router.Group("sysParams").Use(middleware.OperationRecord())
	sysParamsRouterWithoutRecord := Router.Group("sysParams")
	{
		sysParamsRouter.POST("createSysParams", sysParamsApi.CreateSysParams)             // createParameter
		sysParamsRouter.DELETE("deleteSysParams", sysParamsApi.DeleteSysParams)           // delete parameter
		sysParamsRouter.DELETE("deleteSysParamsByIds", sysParamsApi.DeleteSysParamsByIds) // batch delete parameters
		sysParamsRouter.PUT("updateSysParams", sysParamsApi.UpdateSysParams)              // update parameter
	}
	{
		sysParamsRouterWithoutRecord.GET("findSysParams", sysParamsApi.FindSysParams)       // get by IDParameter
		sysParamsRouterWithoutRecord.GET("getSysParamsList", sysParamsApi.GetSysParamsList) // get parameter list
		sysParamsRouterWithoutRecord.GET("getSysParam", sysParamsApi.GetSysParam)           // According toKeygetParameter
	}
}
