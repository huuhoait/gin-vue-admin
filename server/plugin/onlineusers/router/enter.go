package router

import "github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/api"

var (
	Router     = new(routerGroup)
	apiSession = api.Api.Session
)

type routerGroup struct{ Session sessionRouter }
