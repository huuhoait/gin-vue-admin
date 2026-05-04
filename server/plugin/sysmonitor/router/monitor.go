package router

import "github.com/gin-gonic/gin"

type monitorRouter struct{}

func (r *monitorRouter) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	// All read-only — skip OperationRecord to keep poll noise out of audit log.
	group := private.Group("sysmonitor")
	group.GET("server", apiMonitor.Server)
	group.GET("runtime", apiMonitor.Runtime)
	group.GET("cache", apiMonitor.Cache)
	_ = public
}
