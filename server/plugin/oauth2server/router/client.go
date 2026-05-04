package router

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

type clientRouter struct{}

func (r *clientRouter) Init(_ *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("oauth2Client").Use(middleware.OperationRecord())
		group.POST("create", apiClient.CreateClient)
		group.PUT("update", apiClient.UpdateClient)
		group.DELETE("delete", apiClient.DeleteClient)
		group.POST("regenerateSecret", apiClient.RegenerateSecret)
	}
	{
		group := private.Group("oauth2Client")
		group.GET("find", apiClient.FindClient)
		group.GET("list", apiClient.ListClients)
	}
}
