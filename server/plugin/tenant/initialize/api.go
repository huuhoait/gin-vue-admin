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
		{Path: "/tenantMembership/assign", Description: "assign user to tenant", ApiGroup: "TenantMembership", Method: "POST"},
		{Path: "/tenantMembership/unassign", Description: "remove user from tenant", ApiGroup: "TenantMembership", Method: "DELETE"},
		{Path: "/tenantMembership/members", Description: "list members of a tenant", ApiGroup: "TenantMembership", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
