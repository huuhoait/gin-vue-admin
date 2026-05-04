package router

import "github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/api"

var (
	Router        = new(routerGroup)
	apiClient     = api.Api.Client
	apiOAuth2Flow = api.Api.OAuth2
)

type routerGroup struct {
	Client clientRouter
	OAuth2 oauth2Router
}
