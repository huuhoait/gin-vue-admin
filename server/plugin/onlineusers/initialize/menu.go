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
			Path:      "onlineUsers",
			Name:      "onlineUsers",
			Hidden:    false,
			Component: "plugin/onlineusers/view/online.vue",
			Sort:      6,
			Meta:      model.Meta{Title: "admin.plugin.online_users.menu_title", Icon: "user"},
		},
	}
	utils.RegisterMenus(entities...)
}
