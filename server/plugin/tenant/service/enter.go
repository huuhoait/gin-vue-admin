package service

var Service = new(serviceGroup)

type serviceGroup struct {
	Tenant     tenantService
	Membership membershipService
	Package    packageService
}
