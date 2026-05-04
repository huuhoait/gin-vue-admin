package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(_ context.Context) {
	// RegisterMenus treats [0] as parent and [1..] as children. Security groups
	// OAuth2 with Online Sessions (onlineusers plugin registers the same parent).
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
			Path:      "oauth2Clients",
			Name:      "oauth2Clients",
			Hidden:    false,
			Component: "plugin/oauth2server/view/clients.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "admin.plugin.oauth2.menu_title", Icon: "key"},
		},
	}
	utils.RegisterMenus(entities...)
}
