package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/proxy"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	ProxyServiceGroup   *proxy.ServiceGroup
}

// InitProxyServices initialises the proxy clients after config is loaded.
func InitProxyServices() {
	cfg := global.GVA_CONFIG.Proxy
	ServiceGroupApp.ProxyServiceGroup = proxy.NewServiceGroup(
		cfg.CoreServiceURL,
		cfg.OrderServiceURL,
		cfg.RequestTimeout,
	)
}
