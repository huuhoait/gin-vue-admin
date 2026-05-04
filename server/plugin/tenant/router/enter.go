package router

import "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/api"

var (
	Router        = new(routerGroup)
	apiTenant     = api.Api.Tenant
	apiMembership = api.Api.Membership
	apiMy         = api.Api.My
)

type routerGroup struct {
	Tenant     tenantRouter
	Membership membershipRouter
	My         myRouter
}
