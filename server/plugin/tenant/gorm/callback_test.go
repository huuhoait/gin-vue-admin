package gorm

import (
	"context"
	"testing"

	tenantservice "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// scopedRow models a tenant-aware table; embeds a TenantID exactly the way
// production models do via tenantmodel.TenantModel. Defined here so the test
// is self-contained and doesn't depend on plugin migrations.
type scopedRow struct {
	ID       uint   `gorm:"primaryKey"`
	TenantID uint   `gorm:"index;not null"`
	Name     string `gorm:"size:64"`
}

func (scopedRow) TableName() string { return "scoped_rows" }

// globalRow has no TenantID — represents tables outside the tenant boundary
// (e.g. system audit log). The callback must leave these alone.
type globalRow struct {
	ID    uint   `gorm:"primaryKey"`
	Label string `gorm:"size:64"`
}

func (globalRow) TableName() string { return "global_rows" }

// newTestDB builds a fresh in-memory sqlite DB with the tenant callbacks
// installed. Each test gets its own DB to keep state isolated.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&scopedRow{}, &globalRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// Each test installs the callbacks on its own DB; reset the once-guard so
	// Register actually runs.
	resetForTests()
	if err := Register(db); err != nil {
		t.Fatalf("register callbacks: %v", err)
	}
	return db
}

func TestQueryAutoInjectsTenantPredicate(t *testing.T) {
	db := newTestDB(t)

	// Seed two rows in different tenants without scoping (raw insert via a
	// non-tenant context simulates an out-of-band setup).
	seed := []scopedRow{
		{TenantID: 1, Name: "alpha"},
		{TenantID: 2, Name: "bravo"},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	ctx := tenantservice.WithTenant(context.Background(), 1)
	var got []scopedRow
	if err := db.WithContext(ctx).Find(&got).Error; err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got) != 1 || got[0].Name != "alpha" {
		t.Fatalf("expected only tenant 1 rows, got %+v", got)
	}
}

func TestQueryWithoutTenantFieldIsUntouched(t *testing.T) {
	db := newTestDB(t)

	if err := db.Create(&[]globalRow{{Label: "a"}, {Label: "b"}}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	ctx := tenantservice.WithTenant(context.Background(), 7)
	var got []globalRow
	if err := db.WithContext(ctx).Find(&got).Error; err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 rows on a table without TenantID, got %d", len(got))
	}
}

func TestCreateStampsTenantIDFromContext(t *testing.T) {
	db := newTestDB(t)

	ctx := tenantservice.WithTenant(context.Background(), 42)
	row := scopedRow{Name: "stamped"}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if row.TenantID != 42 {
		t.Fatalf("expected stamped TenantID=42, got %d", row.TenantID)
	}

	// Verify the value actually landed in the row, not just the struct: read
	// it back with the bypass marker so we see the raw column.
	bypass := tenantservice.WithTenantIgnore(context.Background())
	var fetched scopedRow
	if err := db.WithContext(bypass).First(&fetched, row.ID).Error; err != nil {
		t.Fatalf("readback: %v", err)
	}
	if fetched.TenantID != 42 {
		t.Fatalf("readback TenantID=%d, want 42", fetched.TenantID)
	}
}

func TestCreateDoesNotOverrideExplicitTenant(t *testing.T) {
	db := newTestDB(t)

	// Caller deliberately wants tenant 9 even though their context says 1
	// (e.g. a backfill script running under an admin context).
	ctx := tenantservice.WithTenant(context.Background(), 1)
	row := scopedRow{TenantID: 9, Name: "explicit"}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if row.TenantID != 9 {
		t.Fatalf("expected explicit TenantID preserved, got %d", row.TenantID)
	}
}

func TestSystemTenantSkipsInjection(t *testing.T) {
	db := newTestDB(t)

	if err := db.Create(&[]scopedRow{{TenantID: 1}, {TenantID: 2}, {TenantID: 3}}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// SystemTenantID (0) means "act as super-admin" — the callback must skip
	// injection so the query sees every tenant's rows.
	ctx := tenantservice.WithTenant(context.Background(), tenantservice.SystemTenantID)
	var got []scopedRow
	if err := db.WithContext(ctx).Find(&got).Error; err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("system tenant should see all rows; got %d", len(got))
	}
}

func TestWithTenantIgnoreSkipsInjection(t *testing.T) {
	db := newTestDB(t)

	if err := db.Create(&[]scopedRow{{TenantID: 5}, {TenantID: 6}}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Caller IS in a tenant context, but explicitly opts out for this query.
	base := tenantservice.WithTenant(context.Background(), 5)
	ctx := tenantservice.WithTenantIgnore(base)

	var got []scopedRow
	if err := db.WithContext(ctx).Find(&got).Error; err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("WithTenantIgnore should bypass; got %d rows", len(got))
	}
}

func TestExplicitTenantPredicateNotOverridden(t *testing.T) {
	db := newTestDB(t)

	if err := db.Create(&[]scopedRow{
		{TenantID: 1, Name: "in-1"},
		{TenantID: 2, Name: "in-2"},
		{TenantID: 3, Name: "in-3"},
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Context says tenant 1, but the caller (e.g. support tool) explicitly
	// queries for tenant 3. The callback must defer to the caller, returning
	// tenant 3's rows — not zero rows from an "AND tenant_id=1" overlay.
	ctx := tenantservice.WithTenant(context.Background(), 1)
	var got []scopedRow
	if err := db.WithContext(ctx).Where("tenant_id = ?", 3).Find(&got).Error; err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got) != 1 || got[0].Name != "in-3" {
		t.Fatalf("expected caller's explicit tenant_id=3 to win, got %+v", got)
	}
}

func TestUpdateAndDeleteAreScoped(t *testing.T) {
	db := newTestDB(t)

	if err := db.Create(&[]scopedRow{
		{TenantID: 1, Name: "original"},
		{TenantID: 2, Name: "other-tenant"},
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	ctx := tenantservice.WithTenant(context.Background(), 1)

	// Update should only touch tenant 1's row even though the WHERE matches
	// the name across tenants.
	if err := db.WithContext(ctx).Model(&scopedRow{}).
		Where("name = ?", "original").
		Update("name", "renamed").Error; err != nil {
		t.Fatalf("update: %v", err)
	}

	bypass := tenantservice.WithTenantIgnore(context.Background())
	var all []scopedRow
	if err := db.WithContext(bypass).Order("tenant_id").Find(&all).Error; err != nil {
		t.Fatalf("readback: %v", err)
	}
	if len(all) != 2 || all[0].Name != "renamed" || all[1].Name != "other-tenant" {
		t.Fatalf("unexpected post-update state: %+v", all)
	}

	// Delete is also scoped: deleting "any row" from tenant 1 must leave
	// tenant 2's row intact.
	if err := db.WithContext(ctx).Where("1 = 1").Delete(&scopedRow{}).Error; err != nil {
		t.Fatalf("delete: %v", err)
	}
	var remaining []scopedRow
	if err := db.WithContext(bypass).Find(&remaining).Error; err != nil {
		t.Fatalf("readback: %v", err)
	}
	if len(remaining) != 1 || remaining[0].TenantID != 2 {
		t.Fatalf("delete should be tenant-scoped; remaining=%+v", remaining)
	}
}

func TestRegisterIsIdempotent(t *testing.T) {
	db := newTestDB(t)
	// Second call must not error or double-register; the once-guard short-
	// circuits even though resetForTests was called inside newTestDB.
	if err := Register(db); err != nil {
		t.Fatalf("second Register call should be a no-op, got %v", err)
	}
}
