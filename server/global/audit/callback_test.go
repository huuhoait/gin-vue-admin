package audit

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type widget struct {
	ID        uint
	Name      string
	CreatedBy uint
	UpdatedBy uint
	DeletedBy uint
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&widget{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	// Reset the once-guard so each test gets a clean callback registration.
	registered.Store(false)
	if err := Register(db); err != nil {
		t.Fatalf("register: %v", err)
	}
	return db
}

func TestStampCreate_FromContext(t *testing.T) {
	db := newDB(t)
	ctx := WithUserID(context.Background(), 42)
	w := widget{Name: "alpha"}
	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	var got widget
	db.First(&got, w.ID)
	if got.CreatedBy != 42 || got.UpdatedBy != 42 {
		t.Fatalf("expected CreatedBy=UpdatedBy=42, got %+v", got)
	}
}

func TestStampCreate_NoContext(t *testing.T) {
	db := newDB(t)
	w := widget{Name: "beta"}
	if err := db.Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if w.CreatedBy != 0 {
		t.Fatalf("expected CreatedBy=0 for unscoped, got %d", w.CreatedBy)
	}
}

func TestStampCreate_PreservesExplicit(t *testing.T) {
	db := newDB(t)
	ctx := WithUserID(context.Background(), 42)
	w := widget{Name: "gamma", CreatedBy: 99} // pretend a backfill
	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	var got widget
	db.First(&got, w.ID)
	if got.CreatedBy != 99 {
		t.Fatalf("explicit CreatedBy=99 should be preserved, got %d", got.CreatedBy)
	}
}

func TestStampUpdate_OverwritesUpdatedBy(t *testing.T) {
	db := newDB(t)
	w := widget{Name: "delta", CreatedBy: 1, UpdatedBy: 1}
	db.Create(&w)

	ctx := WithUserID(context.Background(), 42)
	if err := db.WithContext(ctx).Model(&w).Update("name", "delta-2").Error; err != nil {
		t.Fatalf("update: %v", err)
	}
	var got widget
	db.First(&got, w.ID)
	if got.UpdatedBy != 42 {
		t.Fatalf("expected UpdatedBy=42 after update, got %d", got.UpdatedBy)
	}
	if got.CreatedBy != 1 {
		t.Fatalf("CreatedBy must be preserved on update, got %d", got.CreatedBy)
	}
}

func TestUserIDFromContext_ZeroIgnored(t *testing.T) {
	// WithUserID(ctx, 0) is a no-op so callers don't accidentally clear an
	// inherited userID by passing zero.
	ctx := WithUserID(context.Background(), 0)
	if _, ok := UserIDFromContext(ctx); ok {
		t.Fatal("WithUserID(ctx, 0) should not set the value")
	}

	ctx2 := WithUserID(ctx, 7)
	if id, ok := UserIDFromContext(ctx2); !ok || id != 7 {
		t.Fatalf("expected (7,true), got (%d,%v)", id, ok)
	}
}

func TestRegister_Idempotent(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&widget{})
	registered.Store(false)
	if err := Register(db); err != nil {
		t.Fatal(err)
	}
	// Second call must be a no-op (would otherwise re-register the
	// callback and double-stamp).
	if err := Register(db); err != nil {
		t.Fatal(err)
	}
	ctx := WithUserID(context.Background(), 42)
	w := widget{Name: "epsilon"}
	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatal(err)
	}
	if w.CreatedBy != 42 {
		t.Fatalf("expected CreatedBy=42, got %d", w.CreatedBy)
	}
}
