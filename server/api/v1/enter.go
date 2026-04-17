package v1

import (
	"github.com/huuhoait/gin-vue-admin/server/api/v1/example"
	"github.com/huuhoait/gin-vue-admin/server/api/v1/proxy"
	"github.com/huuhoait/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	ProxyApiGroup   proxy.ApiGroup
}
