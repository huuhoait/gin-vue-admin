# tenant plugin

Multi-tenancy support for the GVA admin: tenant model, user-tenant memberships,
context propagation, and helpers for tenant-aware infrastructure (async work,
Redis caching).

## Layout

- `model/`, `service/`, `api/`, `router/`, `initialize/` — standard GVA layers
  for the tenant + membership CRUD.
- `middleware/context.go` — resolves the active tenant per-request from
  `X-Tenant-ID` (admin override) or the user's primary membership, and stamps
  it onto `c.Request.Context()` via `service.WithTenant`.
- `service/context.go` — `WithTenant` / `FromContext` for context propagation,
  plus `WithTenantScope` / `WithStrictTenantScope` GORM scopes.
- `async/` — propagate tenant context across goroutines and cron callbacks.
- `redis/` — per-tenant Redis key namespacing.

## Async / cron tenant propagation

Goroutines and cron callbacks DO NOT inherit the request's tenant context.
A `go func()` started from a handler runs with a fresh context, so
`service.FromContext(ctx)` returns `(0, false)` and any GORM query that
depends on `WithTenantScope` runs unscoped (system tenant) — best case
data leaks to the wrong tenant, worst case it leaks across tenants.

Use `tenant/async` to capture and propagate explicitly:

```go
import (
    tenantasync "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/async"
)

tenantID := tenantasync.Capture(c.Request.Context())
go tenantasync.Run(tenantID, func(ctx context.Context) {
    // global.GVA_DB queries via service.WithTenantScope(ctx) are now
    // filtered to tenantID, even after the request handler has returned.
})()
```

`Run` deliberately does NOT propagate parent-context cancellation: the goroutine
keeps a fresh background context so the work isn't aborted when the HTTP
handler returns.

For cron tasks pinned to a specific tenant:

```go
global.GVA_Timer.AddTaskByFunc("daily-cleanup", "@daily",
    tenantasync.Run(5, doCleanup),
    "Cleanup tenant 5 stale records",
    cron.WithSeconds(),
)
```

For maintenance jobs that span all tenants, use `tenantasync.AsSystem(fn)` —
it runs with no tenant scope, identical to a super-admin/unscoped query.

```go
global.GVA_Timer.AddTaskByFunc("nightly-rollup", "@daily",
    tenantasync.AsSystem(rollupAllTenants),
    "Cross-tenant nightly rollup",
    cron.WithSeconds(),
)
```

Other helpers:

- `Capture(ctx) uint` — extract the current tenant id (0 if unscoped).
- `RunNow(tenantID, fn)` — synchronous variant; runs `fn` in-line under the
  scoped context.
- `CarryFromContext(ctx, fn)` — shorthand for `Run(Capture(ctx), fn)`.

## Per-tenant Redis namespacing

Redis keys are global by default — `global.GVA_REDIS.Set(ctx, "user:42", …)`
collides between tenants if both have a user with id=42. The DB has the
`tenant_id` discriminator column, but Redis has no such isolation.

Use `tenant/redis.Key(ctx, raw)` to prefix every key with the active tenant
id. System-tenant requests (id=0 / unscoped ctx) get unprefixed keys, so
existing code that hasn't yet been migrated keeps working.

```go
import (
    tenantredis "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/redis"
)

rkey := tenantredis.Key(ctx, "session:" + uuid)            // tenant:5:session:<uuid>
global.GVA_REDIS.Set(ctx, rkey, value, ttl)
```

For SCAN-style operations:

```go
pat := tenantredis.MatchPattern(ctx, "session:*")          // tenant:5:session:*
iter := global.GVA_REDIS.Scan(ctx, 0, pat, 100).Iterator()
for iter.Next(ctx) {
    raw := tenantredis.StripPrefix(ctx, iter.Val())        // session:<uuid>
    // …
}
```

Other helpers:

- `Keys(ctx, raw...)` — batch-prefix a slice; returns a fresh slice.
- `StripPrefix(ctx, key)` — reverse of `Key`; safe no-op when the key is not
  prefixed for the active tenant (never silently rewrites foreign keys).
- `HasPrefix(ctx, key)` — does `key` carry the active tenant's namespace?

The format is fixed (`tenant:<id>:<raw>`) and exposed as the constants
`Prefix` (`"tenant:"`) and `Separator` (`":"`) for tooling that needs to
recognise namespaced entries without re-deriving the format.
