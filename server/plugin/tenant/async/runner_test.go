package async

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
)

// TestCaptureFromScopedContext verifies that Capture extracts the tenant id
// installed by service.WithTenant.
func TestCaptureFromScopedContext(t *testing.T) {
	ctx := service.WithTenant(context.Background(), 7)
	if got := Capture(ctx); got != 7 {
		t.Fatalf("Capture(scoped ctx) = %d, want 7", got)
	}
}

// TestCaptureFromUnscopedContextReturnsZero verifies fail-closed default:
// an unscoped ctx yields tenant id 0 (system tenant).
func TestCaptureFromUnscopedContextReturnsZero(t *testing.T) {
	if got := Capture(context.Background()); got != 0 {
		t.Fatalf("Capture(background) = %d, want 0", got)
	}
}

// TestRunPropagatesTenantToChild verifies that the func returned by Run
// produces a context where service.FromContext sees the requested tenant.
func TestRunPropagatesTenantToChild(t *testing.T) {
	var (
		seenID uint
		seenOK bool
	)
	wrapped := Run(42, func(ctx context.Context) {
		seenID, seenOK = service.FromContext(ctx)
	})
	wrapped()

	if !seenOK {
		t.Fatalf("expected child ctx to be tenant-scoped (ok=true), got ok=false")
	}
	if seenID != 42 {
		t.Fatalf("child saw tenant id %d, want 42", seenID)
	}
}

// TestAsSystemMarksChildAsSystemTenant verifies AsSystem produces a context
// where FromContext returns ok=false (matching the existing
// "id=0 ⇒ system tenant" convention in service.FromContext).
func TestAsSystemMarksChildAsSystemTenant(t *testing.T) {
	var (
		seenID uint
		seenOK bool
	)
	wrapped := AsSystem(func(ctx context.Context) {
		seenID, seenOK = service.FromContext(ctx)
	})
	wrapped()

	// service.FromContext returns ok=false when id==0, which is precisely
	// the system-tenant signal the GORM scope uses to bypass filtering.
	if seenOK {
		t.Fatalf("AsSystem child saw ok=true; system tenant should be unscoped")
	}
	if seenID != 0 {
		t.Fatalf("AsSystem child saw id=%d, want 0", seenID)
	}
}

// TestRunDoesNotPropagateParentCancellation verifies that cancelling the
// parent context does not cancel the child fn. This is the core safety
// guarantee for goroutines spawned from request handlers.
func TestRunDoesNotPropagateParentCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(service.WithTenant(context.Background(), 9))
	tid := Capture(parent)
	cancel() // simulate request handler returning before goroutine completes

	var (
		mu       sync.Mutex
		childErr error
		childID  uint
		childOK  bool
	)

	done := make(chan struct{})
	go Run(tid, func(ctx context.Context) {
		// If the child were attached to the parent's cancellation, ctx.Err()
		// would already be non-nil. It must be nil.
		mu.Lock()
		childErr = ctx.Err()
		childID, childOK = service.FromContext(ctx)
		mu.Unlock()
		close(done)
	})()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("child goroutine did not run within 2s")
	}

	mu.Lock()
	defer mu.Unlock()
	if childErr != nil {
		t.Fatalf("child ctx.Err() = %v, want nil (parent cancellation must NOT propagate)", childErr)
	}
	if !childOK || childID != 9 {
		t.Fatalf("child saw tenant (%d, %v), want (9, true)", childID, childOK)
	}
}

// TestCarryFromContextIsCaptureThenRun verifies the convenience wrapper
// matches Run(Capture(ctx), fn).
func TestCarryFromContextIsCaptureThenRun(t *testing.T) {
	parent := service.WithTenant(context.Background(), 13)

	var seen uint
	CarryFromContext(parent, func(ctx context.Context) {
		seen, _ = service.FromContext(ctx)
	})()

	if seen != 13 {
		t.Fatalf("CarryFromContext propagated id=%d, want 13", seen)
	}
}

// TestRunNowIsSynchronous verifies RunNow executes immediately and applies
// the tenant scoping.
func TestRunNowIsSynchronous(t *testing.T) {
	executed := false
	var seen uint
	RunNow(21, func(ctx context.Context) {
		executed = true
		seen, _ = service.FromContext(ctx)
	})

	if !executed {
		t.Fatal("RunNow did not execute fn synchronously")
	}
	if seen != 21 {
		t.Fatalf("RunNow propagated id=%d, want 21", seen)
	}
}
