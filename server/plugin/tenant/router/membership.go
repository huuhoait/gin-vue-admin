package router

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

type membershipRouter struct{}

func (r *membershipRouter) Init(_ *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("tenantMembership").Use(middleware.OperationRecord())
		group.POST("assign", apiMembership.AssignUser)
		group.DELETE("unassign", apiMembership.UnassignUser)
	}
	{
		group := private.Group("tenantMembership")
		group.GET("members", apiMembership.MembersOfTenant)
	}
}
