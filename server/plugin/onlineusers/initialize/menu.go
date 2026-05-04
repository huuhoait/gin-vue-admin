package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(_ context.Context) {
	// Same Security parent as oauth2server — FirstOrCreate merges by name.
	entities := []model.SysBaseMenu{
		{
			MenuLevel: 0,
			ParentId:  0,
			Path:      "security",
			Name:      "security",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      6,
			Meta:      model.Meta{Title: "admin.menu.security", Icon: "lock"},
		},
		{
			MenuLevel: 1,
			ParentId:  0,
			Path:      "onlineUsers",
			Name:      "onlineUsers",
			Hidden:    false,
			Component: "plugin/onlineusers/view/online.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "admin.plugin.online_users.menu_title", Icon: "user"},
		},
	}
	utils.RegisterMenus(entities...)
}
