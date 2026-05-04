package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(_ context.Context) {
	entities := []model.SysApi{
		{Path: "/sysmonitor/server", Description: "host stats", ApiGroup: "SysMonitor", Method: "GET"},
		{Path: "/sysmonitor/runtime", Description: "go runtime stats", ApiGroup: "SysMonitor", Method: "GET"},
		{Path: "/sysmonitor/cache", Description: "redis info", ApiGroup: "SysMonitor", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
