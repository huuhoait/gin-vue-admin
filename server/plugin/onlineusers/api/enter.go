package api

import "github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/service"

var (
	Api            = new(apiGroup)
	serviceSession = service.Service.Session
)

type apiGroup struct{ Session sessionApi }
