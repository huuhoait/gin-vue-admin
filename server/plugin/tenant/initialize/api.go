package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(_ context.Context) {
	entities := []model.SysApi{
		{Path: "/tenant/create", Description: "create tenant", ApiGroup: "Tenant", Method: "POST"},
		{Path: "/tenant/update", Description: "update tenant", ApiGroup: "Tenant", Method: "PUT"},
		{Path: "/tenant/delete", Description: "delete tenant", ApiGroup: "Tenant", Method: "DELETE"},
		{Path: "/tenant/find", Description: "find tenant", ApiGroup: "Tenant", Method: "GET"},
		{Path: "/tenant/list", Description: "list tenants", ApiGroup: "Tenant", Method: "GET"},
		{Path: "/tenant/mine", Description: "list tenants accessible to the current user", ApiGroup: "Tenant", Method: "GET"},
		{Path: "/tenantMembership/assign", Description: "assign user to tenant", ApiGroup: "TenantMembership", Method: "POST"},
		{Path: "/tenantMembership/unassign", Description: "remove user from tenant", ApiGroup: "TenantMembership", Method: "DELETE"},
		{Path: "/tenantMembership/members", Description: "list members of a tenant", ApiGroup: "TenantMembership", Method: "GET"},
		{Path: "/tenantPackage/create", Description: "create tenant package", ApiGroup: "TenantPackage", Method: "POST"},
		{Path: "/tenantPackage/update", Description: "update tenant package", ApiGroup: "TenantPackage", Method: "PUT"},
		{Path: "/tenantPackage/delete", Description: "delete tenant package", ApiGroup: "TenantPackage", Method: "DELETE"},
		{Path: "/tenantPackage/find", Description: "find tenant package", ApiGroup: "TenantPackage", Method: "GET"},
		{Path: "/tenantPackage/list", Description: "list tenant packages", ApiGroup: "TenantPackage", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
