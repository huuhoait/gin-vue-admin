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
			Path:      "tenants",
			Name:      "tenants",
			Hidden:    false,
			Component: "plugin/tenant/view/tenants.vue",
			Sort:      9,
			Meta:      model.Meta{Title: "admin.plugin.tenant.menu_title", Icon: "office-building"},
		},
	}
	utils.RegisterMenus(entities...)
}
