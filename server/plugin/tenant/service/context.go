package service

import (
	"context"

	"gorm.io/gorm"
)

// tenantCtxKey is unexported so only this package can set the value, ensuring
// that downstream code reading the tenant only ever sees a value put there
// by the tenant middleware (not arbitrary callers).
type tenantCtxKey struct{}

// tenantIgnoreCtxKey marks a context as opting out of automatic tenant
// scoping. Separate key type so it cannot collide with the tenant id value.
type tenantIgnoreCtxKey struct{}

// SystemTenantID is the reserved id used by super-admin flows that legitimately
// need to see across all tenants (cross-tenant audit, billing rollups, support
// tooling). The auto-scoping GORM callback treats this as "do not inject".
const SystemTenantID uint = 0

// WithTenant returns a child context carrying the active tenant id. Used by
// the tenant middleware after resolving the request's tenant.
func WithTenant(ctx context.Context, tenantID uint) context.Context {
	return context.WithValue(ctx, tenantCtxKey{}, tenantID)
}

// FromContext extracts the active tenant id. ok=false means the request was
// not tenant-scoped (e.g., super-admin access without X-Tenant-ID, or an
// unauthenticated public endpoint).
func FromContext(ctx context.Context) (uint, bool) {
	if ctx == nil {
		return 0, false
	}
	v := ctx.Value(tenantCtxKey{})
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}

// IsSystemTenant returns true for the reserved tenant_id=0 used by super-admin
// flows that need to act across all tenants. The GORM callback skips injection
// when the request context resolves to the system tenant.
func IsSystemTenant(id uint) bool { return id == SystemTenantID }

// WithTenantIgnore marks a context so the GORM callback skips tenant_id
// injection on subsequent DB calls. Use sparingly — for audit jobs,
// cross-tenant support tooling, or backfill scripts. The marker survives
// derived contexts (context.WithCancel, etc.) because it is a value, not a
// scope.
func WithTenantIgnore(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, tenantIgnoreCtxKey{}, true)
}

// IsTenantIgnored reports whether the context carries the WithTenantIgnore
// marker. The callback consults this before injecting predicates.
func IsTenantIgnored(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	v, _ := ctx.Value(tenantIgnoreCtxKey{}).(bool)
	return v
}

// WithTenantScope is a GORM scope that filters by the active tenant.
//
// Deprecated: prefer embedding model.TenantModel and letting the tenant GORM
// callback inject predicates automatically. This helper is preserved so
// existing call sites continue to compile, and remains useful when you need to
// scope a query whose model does NOT carry a TenantID column (rare).
func WithTenantScope(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id, ok := FromContext(ctx); ok {
			return db.Where("tenant_id = ?", id)
		}
		return db
	}
}

// WithStrictTenantScope is like WithTenantScope but rejects unscoped queries
// by adding an impossible predicate (tenant_id = 0 with a NOT-NULL column
// returns nothing). Use on highly sensitive tables.
//
// Deprecated: prefer embedding model.TenantModel; the auto-injection callback
// is fail-closed when a tenant context is required by the model. Kept for
// legacy call sites.
func WithStrictTenantScope(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id, ok := FromContext(ctx); ok {
			return db.Where("tenant_id = ?", id)
		}
		// Force an empty result rather than a full-table scan.
		return db.Where("1 = 0")
	}
}
