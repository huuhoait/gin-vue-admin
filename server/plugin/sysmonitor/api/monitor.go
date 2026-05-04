package api

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type monitorApi struct{}

// Server host stats (CPU/RAM/Disk/uptime)
// @Tags     SysMonitor
// @Summary  host stats
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} response.Response{data=service.ServerStats,msg=string}
// @Router   /sysmonitor/server [get]
func (a *monitorApi) Server(c *gin.Context) {
	stats, err := serviceMonitor.Server()
	if err != nil {
		global.GVA_LOG.Error("server stats failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(stats, c)
}

// Runtime Go process stats (goroutines/heap/gc)
// @Tags     SysMonitor
// @Summary  go runtime stats
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} response.Response{data=service.RuntimeStats,msg=string}
// @Router   /sysmonitor/runtime [get]
func (a *monitorApi) Runtime(c *gin.Context) {
	response.OkWithData(serviceMonitor.Runtime(), c)
}

// Cache Redis INFO summary
// @Tags     SysMonitor
// @Summary  redis info summary
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} response.Response{data=service.CacheStats,msg=string}
// @Router   /sysmonitor/cache [get]
func (a *monitorApi) Cache(c *gin.Context) {
	stats, err := serviceMonitor.Cache(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("cache stats failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(stats, c)
}
