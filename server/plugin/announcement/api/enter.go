package api

import "github.com/huuhoait/gin-vue-admin/server/plugin/announcement/service"

var (
	Api         = new(api)
	serviceInfo = service.Service.Info
)

type api struct{ Info info }
