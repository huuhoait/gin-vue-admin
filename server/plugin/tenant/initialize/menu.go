package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(_ context.Context) {
	// RegisterMenus treats menus[0] as a parent and adopts the rest as its
	// children — that's not what we want here since both rows are top-level
	// items under ParentId=9. Two single-element calls keep them as siblings.
	utils.RegisterMenus(model.SysBaseMenu{
		ParentId:  9,
		Path:      "tenants",
		Name:      "tenants",
		Hidden:    false,
		Component: "plugin/tenant/view/tenants.vue",
		Sort:      9,
		Meta:      model.Meta{Title: "admin.plugin.tenant.menu_title", Icon: "office-building"},
	})
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
