package router

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

type packageRouter struct{}

// Init mounts /tenantPackage/* routes. Mutating verbs (create/update/delete)
// are wrapped in OperationRecord middleware; reads are not.
func (r *packageRouter) Init(_ *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("tenantPackage").Use(middleware.OperationRecord())
		group.POST("create", apiTenantPackage.CreatePackage)
		group.PUT("update", apiTenantPackage.UpdatePackage)
		group.DELETE("delete", apiTenantPackage.DeletePackage)
	}
	{
		group := private.Group("tenantPackage")
		group.GET("find", apiTenantPackage.FindPackage)
		group.GET("list", apiTenantPackage.ListPackages)
	}
}
