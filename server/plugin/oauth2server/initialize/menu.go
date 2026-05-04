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
			Path:      "oauth2Clients",
			Name:      "oauth2Clients",
			Hidden:    false,
			Component: "plugin/oauth2server/view/clients.vue",
			Sort:      8,
			Meta:      model.Meta{Title: "admin.plugin.oauth2.menu_title", Icon: "key"},
		},
	}
	utils.RegisterMenus(entities...)
}
