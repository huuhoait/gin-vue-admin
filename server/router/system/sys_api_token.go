package system

import (
	"github.com/huuhoait/gin-vue-admin/server/api/v1"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ApiTokenRouter struct{}

func (s *ApiTokenRouter) InitApiTokenRouter(Router *gin.RouterGroup) {
	apiTokenRouter := Router.Group("sysApiToken").Use(middleware.OperationRecord())
	apiTokenApi := v1.ApiGroupApp.SystemApiGroup.ApiTokenApi
	{
		apiTokenRouter.POST("createApiToken", apiTokenApi.CreateApiToken)   // SignsendToken
		apiTokenRouter.POST("getApiTokenList", apiTokenApi.GetApiTokenList) // list
		apiTokenRouter.POST("deleteApiToken", apiTokenApi.DeleteApiToken)   // VoidToken
	}
}
