package router

import (
	"github.com/huuhoait/gin-vue-admin/server/router/example"
	"github.com/huuhoait/gin-vue-admin/server/router/proxy"
	"github.com/huuhoait/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Proxy   proxy.RouterGroup
}
