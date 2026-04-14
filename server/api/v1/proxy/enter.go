package proxy

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/service/dashboard"
	proxyPkg "github.com/flipped-aurora/gin-vue-admin/server/service/proxy"
	"gorm.io/gorm"
)

type ApiGroup struct {
	SkyAgentApi
	DashboardApi
}

// getProxyService returns the lazily-initialised proxy service group.
func getProxyService() *proxyPkg.ServiceGroup {
	if service.ServiceGroupApp.ProxyServiceGroup == nil {
		service.InitProxyServices()
	}
	return service.ServiceGroupApp.ProxyServiceGroup
}

// --- Dashboard service wiring (avoids import cycle with initialize) ---

// DashboardDBs is set by initialize package after DB connections are ready.
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
