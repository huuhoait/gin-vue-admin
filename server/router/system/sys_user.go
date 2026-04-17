package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("user")
	{
		userRouter.POST("admin_register", baseApi.Register)               // managementMemberRegisterAccount
		userRouter.POST("changePassword", baseApi.ChangePassword)         // user modifiedpassword
		userRouter.POST("setUserAuthority", baseApi.SetUserAuthority)     // set user authority
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)               // delete user
		userRouter.PUT("setUserInfo", baseApi.SetUserInfo)                // set user info
		userRouter.PUT("setSelfInfo", baseApi.SetSelfInfo)                // setown profile
		userRouter.POST("setUserAuthorities", baseApi.SetUserAuthorities) // set user authorityGroup
		userRouter.POST("resetPassword", baseApi.ResetPassword)           // reset user password
		userRouter.PUT("setSelfSetting", baseApi.SetSelfSetting)          // user interfaceconfiguration
	}
	{
		userRouterWithoutRecord.POST("getUserList", baseApi.GetUserList) // Paginationget user list
		userRouterWithoutRecord.GET("getUserInfo", baseApi.GetUserInfo)  // getown profile
	}
}
