package service

var Service = new(serviceGroup)

type serviceGroup struct {
	Client clientService
	OAuth2 oauth2Service
}
