package redis

import (
	"context"
	"reflect"
	"testing"

	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
)

// TestKeyOnUnscopedContextReturnsRaw verifies that an unscoped ctx
// (no tenant set) leaves keys unprefixed for backwards-compat with
// pre-tenancy code paths.
func TestKeyOnUnscopedContextReturnsRaw(t *testing.T) {
	got := Key(context.Background(), "user:42")
	if got != "user:42" {
		t.Fatalf("Key(unscoped, %q) = %q, want unchanged", "user:42", got)
	}
}

// TestKeyOnSystemTenantContextReturnsRaw verifies that ctx scoped to
// tenant id 0 (system tenant) is treated identically to unscoped.
func TestKeyOnSystemTenantContextReturnsRaw(t *testing.T) {
	// service.WithTenant(_, 0) is allowed but service.FromContext returns
	// ok=false for id==0 — i.e. the system tenant signal. Keys must NOT
	// be prefixed in this case.
	ctx := service.WithTenant(context.Background(), 0)
	got := Key(ctx, "user:42")
	if got != "user:42" {
		t.Fatalf("Key(system, %q) = %q, want unchanged", "user:42", got)
	}
}

// TestKeyOnTenantContextPrefixes verifies the basic happy path.
func TestKeyOnTenantContextPrefixes(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 5)
	got := Key(ctx, "user:42")
	if got != "tenant:5:user:42" {
		t.Fatalf("Key(tenant=5, %q) = %q, want %q", "user:42", got, "tenant:5:user:42")
	}
}

// TestKeysBatchPrefixes verifies Keys applies the prefix to every entry
// and returns a fresh slice (defensive copy semantics).
func TestKeysBatchPrefixes(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 3)
	in := []string{"a", "b", "c"}
	got := Keys(ctx, in...)
	want := []string{"tenant:3:a", "tenant:3:b", "tenant:3:c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys(tenant=3, %v) = %v, want %v", in, got, want)
	}
	// Mutating the returned slice must not touch the input.
	got[0] = "mutated"
	if in[0] != "a" {
		t.Fatalf("Keys mutated input slice: in[0]=%q", in[0])
	}
}

// TestKeysBatchOnUnscopedContextReturnsCopy verifies Keys still returns a
// fresh slice when no prefixing happens.
func TestKeysBatchOnUnscopedContextReturnsCopy(t *testing.T) {
	in := []string{"a", "b"}
	got := Keys(context.Background(), in...)
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("Keys(unscoped) = %v, want %v", got, in)
	}
	got[0] = "mutated"
	if in[0] != "a" {
		t.Fatalf("Keys mutated input slice: in[0]=%q", in[0])
	}
}

// TestStripPrefixReverseOfKey verifies StripPrefix removes the namespace
// when the key matches the active tenant's prefix.
func TestStripPrefixReverseOfKey(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 11)
	full := Key(ctx, "session:abc")
	got := StripPrefix(ctx, full)
	if got != "session:abc" {
		t.Fatalf("StripPrefix(%q) = %q, want %q", full, got, "session:abc")
	}
}

// TestStripPrefixUnscopedIsNoop verifies the no-op behaviour on unscoped
// and system-tenant contexts.
func TestStripPrefixUnscopedIsNoop(t *testing.T) {
	if got := StripPrefix(context.Background(), "tenant:5:foo"); got != "tenant:5:foo" {
		t.Fatalf("StripPrefix(unscoped) = %q, want pass-through", got)
	}
	sysCtx := service.WithTenant(context.Background(), 0)
	if got := StripPrefix(sysCtx, "tenant:5:foo"); got != "tenant:5:foo" {
		t.Fatalf("StripPrefix(system) = %q, want pass-through", got)
	}
}

// TestStripPrefixForeignTenantIsNoop verifies that StripPrefix never
// silently rewrites a key from a different tenant — the input string
// must be returned unchanged when its prefix doesn't match the active
// tenant. Otherwise StripPrefix could quietly hand callers data from
// another tenant when, e.g., a misconfigured SCAN spans tenants.
func TestStripPrefixForeignTenantIsNoop(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 7)
	foreign := "tenant:9:session:abc"
	if got := StripPrefix(ctx, foreign); got != foreign {
		t.Fatalf("StripPrefix(tenant=7, foreign=%q) = %q, want unchanged", foreign, got)
	}
}

// TestMatchPatternWrapsKey verifies MatchPattern is consistent with Key
// (it is intentionally just a documenting alias).
func TestMatchPatternWrapsKey(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 5)
	if got, want := MatchPattern(ctx, "session:*"), "tenant:5:session:*"; got != want {
		t.Fatalf("MatchPattern(tenant=5, %q) = %q, want %q", "session:*", got, want)
	}
	// Unscoped ctx returns the raw pattern.
	if got, want := MatchPattern(context.Background(), "session:*"), "session:*"; got != want {
		t.Fatalf("MatchPattern(unscoped, %q) = %q, want %q", "session:*", got, want)
	}
}

// TestHasPrefixDetectsActiveTenantOnly verifies HasPrefix returns true
// only when key is namespaced for the active tenant.
func TestHasPrefixDetectsActiveTenantOnly(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 4)
	if !HasPrefix(ctx, "tenant:4:foo") {
		t.Fatalf("HasPrefix should detect own-tenant prefix")
	}
	if HasPrefix(ctx, "tenant:5:foo") {
		t.Fatalf("HasPrefix must not match foreign-tenant prefix")
	}
	if HasPrefix(ctx, "foo") {
		t.Fatalf("HasPrefix must not match unprefixed key")
	}
	// System/unscoped ctx never reports HasPrefix=true (system keys are
	// not prefixed; the question is meaningless).
	if HasPrefix(context.Background(), "tenant:4:foo") {
		t.Fatalf("HasPrefix on unscoped ctx should be false")
	}
}
