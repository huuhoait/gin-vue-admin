package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	coremw "github.com/huuhoait/gin-vue-admin/server/middleware"
	tenantmw "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/middleware"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/router"

	"github.com/gin-gonic/gin"
)

// Router wires the tenant admin endpoints. Note that the TenantContext
// middleware is applied here; downstream plugins that want their own routes
// to be tenant-scoped should mount tenantmw.TenantContext() on their
// PrivateGroup themselves (or the framework can register it globally — see
// plugin.go init).
func Router(engine *gin.Engine) {
	public := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	private.Use(coremw.JWTAuth()).Use(coremw.CasbinHandler()).Use(tenantmw.TenantContext())
	router.Router.Tenant.Init(public, private)
	router.Router.Membership.Init(public, private)
	router.Router.My.Init(public, private)
}
