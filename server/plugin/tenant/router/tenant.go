package router

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

type tenantRouter struct{}

func (r *tenantRouter) Init(_ *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("tenant").Use(middleware.OperationRecord())
		group.POST("create", apiTenant.CreateTenant)
		group.PUT("update", apiTenant.UpdateTenant)
		group.DELETE("delete", apiTenant.DeleteTenant)
	}
	{
		group := private.Group("tenant")
		group.GET("find", apiTenant.FindTenant)
		group.GET("list", apiTenant.ListTenants)
	}
}
