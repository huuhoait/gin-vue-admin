package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	apiRouter := Router.Group("api").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("api")

	apiPublicRouterWithoutRecord := RouterPub.Group("api")
	{
		apiRouter.GET("getApiGroups", apiRouterApi.GetApiGroups)          // getroute group
		apiRouter.GET("syncApi", apiRouterApi.SyncApi)                    // SynchronousApi
		apiRouter.POST("ignoreApi", apiRouterApi.IgnoreApi)               // IgnoreApi
		apiRouter.POST("enterSyncApi", apiRouterApi.EnterSyncApi)         // sureAcknowledgeSynchronousApi
		apiRouter.POST("createApi", apiRouterApi.CreateApi)               // CreateApi
		apiRouter.POST("deleteApi", apiRouterApi.DeleteApi)               // deleteApi
		apiRouter.POST("getApiById", apiRouterApi.GetApiById)             // getDocumentRowApiMessage
		apiRouter.POST("updateApi", apiRouterApi.UpdateApi)               // update API
		apiRouter.DELETE("deleteApisByIds", apiRouterApi.DeleteApisByIds) // delete selected APIs
		apiRouter.POST("setApiRoles", apiRouterApi.SetApiRoles)          // full overwriteAPIAssociationRole
	}
	{
		apiRouterWithoutRecord.POST("getAllApis", apiRouterApi.GetAllApis) // get all APIs
		apiRouterWithoutRecord.POST("getApiList", apiRouterApi.GetApiList) // getApiList
		apiRouterWithoutRecord.GET("getApiRoles", apiRouterApi.GetApiRoles) // getAPIAssociationrole IDList
	}
	{
		apiPublicRouterWithoutRecord.GET("freshCasbin", apiRouterApi.FreshCasbin) // RefreshcasbinPermission
	}
}
