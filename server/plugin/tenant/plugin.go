// Package tenant provides multi-tenant scaffolding as an opt-in plugin.
//
// Adoption pattern for downstream models:
//
//  1. Embed model.TenantModel instead of global.GVA_MODEL on the rows that
//     belong to a single tenant. The mixin contributes a TenantID column.
//  2. Apply middleware.TenantContext() on private route groups so the active
//     tenant is attached to the request context (the tenant plugin's own
//     routes already do this).
//  3. After that, queries / updates / deletes / creates against tenant-aware
//     models are automatically scoped by the GORM callback installed by this
//     plugin — no per-call WithTenantScope needed. Cross-tenant access for
//     super-admin flows is enabled by either the tenant_id=0 system tenant
//     or service.WithTenantIgnore(ctx) for one-off bypasses.
//
// The plugin does not modify gva_users; tenant ↔ user links live in
// gva_user_tenants and are resolved per-request.
package tenant

import (
	"context"

	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/initialize"
	interfaces "github.com/huuhoait/gin-vue-admin/server/utils/plugin/v2"

	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() { interfaces.Register(Plugin) }

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	initialize.Api(ctx)
	initialize.Menu(ctx)
	// Gorm() runs AutoMigrate AND wires the auto-scoping GORM callbacks. It
	// must run before any tenant-aware queries fly, which is the case here:
	// Register() is called during framework boot, before route handlers serve
	// traffic.
	initialize.Gorm(ctx)
	initialize.Router(group)
}
