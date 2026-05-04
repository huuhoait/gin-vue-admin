package router

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type sessionRouter struct{}

func (r *sessionRouter) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		// Mutating: kick is recorded in the operation log.
		group := private.Group("onlineUsers").Use(middleware.OperationRecord())
		group.POST("kick", apiSession.KickSession)
	}
	{
		// Read-only: skip OperationRecord to avoid log noise from polling.
		group := private.Group("onlineUsers")
		group.GET("list", apiSession.ListSessions)
	}
	_ = public
}
