package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/huuhoait/gin-vue-admin/server/plugin/skyagent/router"

	"github.com/gin-gonic/gin"
)

// Router mounts the SkyAgent BFF endpoints under
// {RouterPrefix}/admin-api/v1 with JWT + Casbin protection. The mount
// path mirrors the FE↔BFF contract in external-frontend-integration.md.
func Router(engine *gin.Engine) {
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix + "/admin-api/v1")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	router.Router.SkyAgent.Init(private)
}
