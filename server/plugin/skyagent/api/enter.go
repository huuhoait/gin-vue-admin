package api

import (
	"sync"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/skyagent/dashboard"
	proxyPkg "github.com/huuhoait/gin-vue-admin/server/plugin/skyagent/service"

	"gorm.io/gorm"
)

// Api is the plugin-local singleton wired into router.Router.SkyAgent. The
// shape mirrors GVA's standard ApiGroup pattern; keeping it inside the plugin
// keeps the cross-plugin coupling at zero.
var Api = struct {
	SkyAgent  SkyAgentApi
	Dashboard DashboardApi
}{}

// --- proxy service wiring ----------------------------------------------------

// proxyServices is the lazily-initialised proxy client group. Built from
// global.GVA_CONFIG.Proxy on first access so plugin import order doesn't have
// to wait for viper.
var (
	proxyServices     *proxyPkg.ServiceGroup
	proxyServicesOnce sync.Once
)

// getProxyService returns the lazily-initialised proxy service group.
func getProxyService() *proxyPkg.ServiceGroup {
	proxyServicesOnce.Do(func() {
		cfg := global.GVA_CONFIG.Proxy
		proxyServices = proxyPkg.NewServiceGroup(
			cfg.CoreServiceURL,
			cfg.OrderServiceURL,
			cfg.RequestTimeout,
		)
	})
	return proxyServices
}

// --- Dashboard service wiring (avoids import cycle with initialize) ---------

// DashboardDBs is set by the plugin's initialize.Readonly after the readonly
// connections are opened. Exported so initialize can populate it without
// re-importing this package's internals.
var DashboardDBs struct {
	CoreDB  *gorm.DB
	OrderDB *gorm.DB
}

var dashboardSvc *dashboard.Service

func getDashboardService() *dashboard.Service {
	if dashboardSvc != nil {
		return dashboardSvc
	}
	repo := dashboard.NewRepository(DashboardDBs.CoreDB, DashboardDBs.OrderDB)
	cache := dashboard.NewCache(nil)
	dashboardSvc = dashboard.NewService(repo, cache)
	return dashboardSvc
}
