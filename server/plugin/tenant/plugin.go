// Package tenant provides multi-tenant scaffolding as an opt-in plugin.
//
// Adoption pattern for downstream models:
//
//  1. Embed model.TenantModel instead of global.GVA_MODEL on the rows that
//     belong to a single tenant.
//  2. In service queries, attach service.WithTenantScope(ctx) (or
//     WithStrictTenantScope for sensitive tables) so reads/updates are
//     auto-filtered.
//  3. On Create, populate TenantID from service.FromContext(ctx) — the
//     plugin does not auto-stamp create paths because it cannot know which
//     create operations belong to a tenant.
//  4. Apply middleware.TenantContext() on private route groups that need
//     scoping (the tenant plugin's own routes already do this).
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
	initialize.Gorm(ctx)
	initialize.Router(group)
}
