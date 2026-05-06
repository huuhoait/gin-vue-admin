package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(_ context.Context) {
	// Tenants is a parent menu. Members management becomes a submenu so primary
	// tenant members can manage their current tenant without opening a drawer.
	utils.RegisterMenus(
		model.SysBaseMenu{
			ParentId:  9,
			Path:      "tenants",
			Name:      "tenants",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      9,
			Meta:      model.Meta{Title: "admin.plugin.tenant.menu_title", Icon: "office-building"},
		},
		model.SysBaseMenu{
			Path:      "list",
			Name:      "tenantList",
			Hidden:    false,
			Component: "plugin/tenant/view/tenants.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "admin.plugin.tenant.menu_title", Icon: "office-building"},
		},
		model.SysBaseMenu{
			Path:      "members",
			Name:      "tenantMembers",
			Hidden:    false,
			Component: "plugin/tenant/view/members.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "admin.plugin.tenant.menu_members_title", Icon: "user"},
		},
	)
	utils.RegisterMenus(model.SysBaseMenu{
		ParentId:  9,
		Path:      "tenantPackages",
		Name:      "tenantPackages",
		Hidden:    false,
		Component: "plugin/tenant/view/packages.vue",
		Sort:      10,
		Meta:      model.Meta{Title: "admin.plugin.tenant.menu_packages_title", Icon: "files"},
	})
}
