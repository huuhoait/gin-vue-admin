// Package async provides helpers that propagate tenant context across
// goroutines and cron callbacks. The default behaviour is to FAIL CLOSED:
// goroutines spawned via Run that don't explicitly carry a tenant run as
// the system tenant (id=0), preserving the same default-to-system-tenant
// semantics the GORM tenant scope uses for direct request handlers — i.e.
// service.FromContext returns ok=false and queries are unscoped.
//
// Why a fresh background context?
//
// Request-bound contexts (gin.Context, *http.Request.Context) are cancelled
// the moment the HTTP handler returns. If a goroutine simply captured the
// parent ctx, the in-flight DB write or downstream call would be cut
// mid-stream. Run therefore detaches lifecycle from the parent (using
// context.Background) while preserving the tenant scoping value.
//
// Typical usage:
//
//	tenantID := async.Capture(c.Request.Context())
//	go async.Run(tenantID, func(ctx context.Context) {
//	    // global.GVA_DB queries that go through service.WithTenantScope(ctx)
//	    // will be filtered to tenantID — even though c.Request is gone.
//	})()
package async

import (
	"context"

	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
)

// Capture extracts the current tenant from ctx, returning 0 (system tenant)
// when ctx is unscoped. Use right before launching a goroutine so the
// parent ctx's cancellation is decoupled from the child's tenant scoping.
//
// Capture intentionally squashes the (uint, bool) shape of FromContext to
// a single uint: callers that need to spawn work always need *some*
// tenant-id to bind, and the system-tenant default (0) is the same
// fail-closed scoping the GORM helpers already use for unscoped requests.
func Capture(ctx context.Context) uint {
	id, ok := service.FromContext(ctx)
	if !ok {
		return 0
	}
	return id
}

// Run executes fn in a fresh background context scoped to tenantID. It
// returns a func() with no arguments so the result can be passed straight
// to `go` or to cron registration:
//
//	tenantID := async.Capture(c.Request.Context())
//	go async.Run(tenantID, func(ctx context.Context) {
//	    // queries through global.GVA_DB scoped to tenantID
//	})()
//
// Pass tenantID=0 to opt into system-tenant scope (cross-tenant
// maintenance jobs). Prefer AsSystem in that case so call sites read
// clearly.
func Run(tenantID uint, fn func(context.Context)) func() {
	return func() {
		ctx := service.WithTenant(context.Background(), tenantID)
		fn(ctx)
	}
}

// RunNow is the synchronous variant — invokes fn immediately with the
// scoped context rather than returning a deferred func(). It is convenient
// when the caller wants to execute work in-line under a different tenant
// scope, e.g. inside a maintenance batch loop.
//
//	for _, tid := range allTenants {
//	    async.RunNow(tid, func(ctx context.Context) { reindex(ctx) })
//	}
func RunNow(tenantID uint, fn func(context.Context)) {
	ctx := service.WithTenant(context.Background(), tenantID)
	fn(ctx)
}

// CarryFromContext is a convenience: capture from ctx, hand back a wrapper
// suitable for `go` or cron. Equivalent to Run(Capture(ctx), fn).
//
//	go async.CarryFromContext(c.Request.Context(), func(ctx context.Context) {
//	    publishOperationLog(ctx, entry)
//	})()
func CarryFromContext(ctx context.Context, fn func(context.Context)) func() {
	return Run(Capture(ctx), fn)
}

// AsSystem runs fn as the system tenant (cross-tenant scope). Equivalent
// to Run(0, fn). Spelled out so call sites read clearly at the start of
// admin/maintenance jobs that intentionally span tenants.
//
//	global.GVA_Timer.AddTaskByFunc("nightly-rollup", "0 0 0 * * *",
//	    async.AsSystem(rollupAllTenants),
//	    "Cross-tenant nightly rollup",
//	    cron.WithSeconds(),
//	)
func AsSystem(fn func(context.Context)) func() {
	return Run(0, fn)
}
