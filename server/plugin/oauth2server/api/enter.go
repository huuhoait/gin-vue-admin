package api

import "github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/service"

var (
	Api            = new(apiGroup)
	serviceClient  = service.Service.Client
	serviceOAuth2  = service.Service.OAuth2
)

type apiGroup struct {
	Client clientApi
	OAuth2 oauth2Api
}
