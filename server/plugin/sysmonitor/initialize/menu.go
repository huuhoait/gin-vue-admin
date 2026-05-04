package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(_ context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  9,
			Path:      "sysmonitor",
			Name:      "sysmonitor",
			Hidden:    false,
			Component: "plugin/sysmonitor/view/dashboard.vue",
			Sort:      7,
			Meta:      model.Meta{Title: "admin.plugin.sysmonitor.menu_title", Icon: "monitor"},
		},
	}
	utils.RegisterMenus(entities...)
}
