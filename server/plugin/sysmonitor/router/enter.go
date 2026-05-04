package router

import "github.com/huuhoait/gin-vue-admin/server/plugin/sysmonitor/api"

var (
	Router     = new(routerGroup)
	apiMonitor = api.Api.Monitor
)

type routerGroup struct{ Monitor monitorRouter }
