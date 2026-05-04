package api

import "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"

var (
	Api                  = new(apiGroup)
	serviceTenant        = service.Service.Tenant
	serviceMembership    = service.Service.Membership
	serviceTenantPackage = service.Service.Package
)

type apiGroup struct {
	Tenant     tenantApi
	Membership membershipApi
	Package    packageApi
	My         myApi
}
