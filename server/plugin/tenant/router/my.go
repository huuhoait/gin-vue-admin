package router

import "github.com/gin-gonic/gin"

type myRouter struct{}

// Init wires the per-user tenant lookup endpoint. It deliberately skips the
// OperationRecord middleware because it is a high-frequency read used by the
// frontend tenant switcher and would pollute the audit log.
func (r *myRouter) Init(_ *gin.RouterGroup, private *gin.RouterGroup) {
	group := private.Group("tenant")
	group.GET("mine", apiMy.MyTenants)
}
