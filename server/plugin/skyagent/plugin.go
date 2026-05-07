// Package skyagent is the BFF plugin for SkyAgent Core/Order. It owns the
// `/admin-api/v1/*` proxy surface, the dashboard read-only DBs, and the
// agent/order/catalog/onboarding admin views consumed by the FE plugin at
// web/src/plugin/skyagent.
//
// Config still lives in core (server/config/proxy.go) — these are URLs +
// timeouts, not plugin-private secrets, and the dashboard read-only DSNs are
// optional.
package skyagent

import (
	"github.com/huuhoait/gin-vue-admin/server/plugin/skyagent/initialize"
	interfaces "github.com/huuhoait/gin-vue-admin/server/utils/plugin/v2"

	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() { interfaces.Register(Plugin) }

func (p *plugin) Register(engine *gin.Engine) {
	// Open the dashboard's read-only Core/Order Postgres connections, if
	// configured. Best-effort; the dashboard endpoint degrades gracefully when
	// the DSNs are blank or unreachable.
	initialize.Readonly()
	initialize.Router(engine)
}
