package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AuthorityRouter struct{}

func (s *AuthorityRouter) InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	authorityRouterWithoutRecord := Router.Group("authority")
	{
		authorityRouter.POST("createAuthority", authorityApi.CreateAuthority)   // create role
		authorityRouter.POST("deleteAuthority", authorityApi.DeleteAuthority)   // delete role
		authorityRouter.PUT("updateAuthority", authorityApi.UpdateAuthority)    // updateRole
		authorityRouter.POST("copyAuthority", authorityApi.CopyAuthority)       // copy role
		authorityRouter.POST("setDataAuthority", authorityApi.SetDataAuthority) // set role resource permissions
		authorityRouter.POST("setRoleUsers", authorityApi.SetRoleUsers)         // full overwriteRoleassociated user
	}
	{
		authorityRouterWithoutRecord.POST("getAuthorityList", authorityApi.GetAuthorityList)     // get role list
		authorityRouterWithoutRecord.GET("getUsersByAuthority", authorityApi.GetUsersByAuthority) // getRoleassociated userIDList
	}
}
