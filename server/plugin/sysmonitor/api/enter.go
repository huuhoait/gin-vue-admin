package api

import "github.com/huuhoait/gin-vue-admin/server/plugin/sysmonitor/service"

var (
	Api            = new(apiGroup)
	serviceMonitor = service.Service.Monitor
)

type apiGroup struct{ Monitor monitorApi }
