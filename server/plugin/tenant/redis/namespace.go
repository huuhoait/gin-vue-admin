// Package redis provides per-tenant Redis key namespacing. It does NOT wrap
// the redis.UniversalClient — callers continue to use global.GVA_REDIS
// directly. They just route their key strings through Key/Keys before
// every read/write, so cache entries from different tenants cannot collide
// on a shared Redis instance.
//
// Keys for the system tenant (id=0, i.e. an unscoped ctx) are NOT prefixed,
// preserving backwards compatibility for code paths that don't yet
// propagate tenant context (e.g. boot-time caches, framework-internal
// keys, anything intentionally scoped at the platform level).
//
// Format:  tenant:<id>:<raw>
//
// Examples:
//
//	rkey := tredis.Key(ctx, "session:" + uuid)
//	global.GVA_REDIS.Set(ctx, rkey, value, ttl)
//
//	pat := tredis.MatchPattern(ctx, "session:*")
//	iter := global.GVA_REDIS.Scan(ctx, 0, pat, 100).Iterator()
package redis

import (
	"context"
	"strconv"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
)

// Prefix is the literal namespace marker prepended to keys. Exposed so
// tooling (e.g. cache-warmers, debug dumps) can recognise namespaced
// entries without re-deriving the format.
const Prefix = "tenant:"

// Separator joins the tenant id segment to the raw key.
const Separator = ":"

// activeTenant returns the tenant id stamped on ctx, or 0 when ctx is
// unscoped or scoped to the system tenant. Both signals collapse to the
// same "do not prefix" decision.
func activeTenant(ctx context.Context) uint {
	id, ok := service.FromContext(ctx)
	if !ok {
		return 0
	}
	return id
}

// prefixFor builds the per-tenant prefix segment "tenant:<id>:" used by
// Key/Keys/MatchPattern. Cached locally per call — strconv allocations
// here are negligible vs the Redis round-trip the prefix prepares for.
func prefixFor(id uint) string {
	return Prefix + strconv.FormatUint(uint64(id), 10) + Separator
}

// Key returns "tenant:<id>:<raw>" when ctx is scoped to a non-system
// tenant, otherwise returns raw unchanged.
//
// Use this on every Redis key written or read from a tenant-aware code
// path. It is a pure string transform; cheap enough to call inline at
// every call-site.
func Key(ctx context.Context, raw string) string {
	id := activeTenant(ctx)
	if id == 0 {
		return raw
	}
	return prefixFor(id) + raw
}

// Keys batch-prefixes a slice of raw keys. Returns a new slice so the
// caller's input is not mutated. When ctx is unscoped/system, the result
// is a defensive copy of raw.
func Keys(ctx context.Context, raw ...string) []string {
	out := make([]string, len(raw))
	id := activeTenant(ctx)
	if id == 0 {
		copy(out, raw)
		return out
	}
	p := prefixFor(id)
	for i, k := range raw {
		out[i] = p + k
	}
	return out
}

// StripPrefix is the reverse of Key — given a namespaced key, return the
// original raw key. Useful for SCAN result processing where the caller
// knows they only scanned within their own tenant and wants to surface
// the un-namespaced form back to business code.
//
// On unscoped/system ctx StripPrefix is a no-op. If the supplied key does
// not begin with the active tenant's prefix, it is returned unchanged
// (StripPrefix never silently rewrites keys from other tenants).
func StripPrefix(ctx context.Context, key string) string {
	id := activeTenant(ctx)
	if id == 0 {
		return key
	}
	return strings.TrimPrefix(key, prefixFor(id))
}

// MatchPattern produces a SCAN/KEYS pattern scoped to the active tenant.
// E.g. for raw "session:*" returns "tenant:5:session:*". This is just
// Key under a different name, but it documents intent at call sites that
// build glob patterns for SCAN.
func MatchPattern(ctx context.Context, raw string) string {
	return Key(ctx, raw)
}

// HasPrefix reports whether key is namespaced for the active tenant on
// ctx. False for system-tenant ctx (where no prefixing is performed) and
// false for keys scoped to a different tenant.
func HasPrefix(ctx context.Context, key string) bool {
	id := activeTenant(ctx)
	if id == 0 {
		return false
	}
	return strings.HasPrefix(key, prefixFor(id))
}
