package service

import (
	"context"

	"gorm.io/gorm"
)

// tenantCtxKey is unexported so only this package can set the value, ensuring
// that downstream code reading the tenant only ever sees a value put there
// by the tenant middleware (not arbitrary callers).
type tenantCtxKey struct{}

// WithTenant returns a child context carrying the active tenant id. Used by
// the tenant middleware after resolving the request's tenant.
func WithTenant(ctx context.Context, tenantID uint) context.Context {
	return context.WithValue(ctx, tenantCtxKey{}, tenantID)
}

// FromContext extracts the active tenant id. ok=false means the request was
// not tenant-scoped (e.g., super-admin access without X-Tenant-ID, or an
// unauthenticated public endpoint).
func FromContext(ctx context.Context) (uint, bool) {
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

// WithTenantScope is a GORM scope that filters by the active tenant. It is a
// no-op when the context is not tenant-scoped — that lets super-admin
// queries see all rows. To force strict isolation (no row leakage even for
// super-admin), use WithStrictTenantScope below.
//
// Usage:
//
//	db := global.GVA_DB.Model(&Order{}).Scopes(service.WithTenantScope(ctx))
//	db.Find(&orders)
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
func WithStrictTenantScope(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id, ok := FromContext(ctx); ok {
			return db.Where("tenant_id = ?", id)
		}
		// Force an empty result rather than a full-table scan.
		return db.Where("1 = 0")
	}
}
