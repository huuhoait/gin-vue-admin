package service

import (
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
)

// setupTestDB swaps global.GVA_DB for an isolated in-memory sqlite handle so
// each test gets a fresh schema. The previous DB is restored on cleanup so
// suites that share the global don't leak state into each other.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_pragma=foreign_keys(1)"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.Tenant{}, &model.UserTenant{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	prev := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		global.GVA_DB = prev
		// Drop tables so a re-used cache:shared sqlite handle starts clean.
		_ = db.Migrator().DropTable(&model.UserTenant{}, &model.Tenant{})
	})
	return db
}

func TestIsActive_DisabledTenant(t *testing.T) {
	setupTestDB(t)
	svc := tenantService{}
	tenant := model.Tenant{Enabled: false}
	if svc.IsActive(tenant) {
		t.Fatalf("expected disabled tenant to be inactive")
	}
}

func TestIsActive_ExpiredTenant(t *testing.T) {
	setupTestDB(t)
	svc := tenantService{}
	past := time.Now().Add(-1 * time.Hour)
	tenant := model.Tenant{Enabled: true, ExpireAt: &past}
	if svc.IsActive(tenant) {
		t.Fatalf("expected expired tenant to be inactive")
	}
}

func TestIsActive_NeverExpires(t *testing.T) {
	setupTestDB(t)
	svc := tenantService{}
	tenant := model.Tenant{Enabled: true, ExpireAt: nil}
	if !svc.IsActive(tenant) {
		t.Fatalf("expected enabled, no-expiration tenant to be active")
	}
}

func TestIsActive_FutureExpiration(t *testing.T) {
	setupTestDB(t)
	svc := tenantService{}
	future := time.Now().Add(24 * time.Hour)
	tenant := model.Tenant{Enabled: true, ExpireAt: &future}
	if !svc.IsActive(tenant) {
		t.Fatalf("expected tenant with future expiration to be active")
	}
}

func TestAssign_AccountLimitReached(t *testing.T) {
	setupTestDB(t)
	tSvc := tenantService{}
	mSvc := membershipService{}

	tenant, err := tSvc.Create(request.CreateTenantReq{
		Code:         "acme",
		Name:         "Acme",
		AccountLimit: 2,
	})
	if err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	if err := mSvc.Assign(101, tenant.ID, true); err != nil {
		t.Fatalf("first assign: %v", err)
	}
	if err := mSvc.Assign(102, tenant.ID, false); err != nil {
		t.Fatalf("second assign: %v", err)
	}

	err = mSvc.Assign(103, tenant.ID, false)
	if !errors.Is(err, ErrAccountLimitReached) {
		t.Fatalf("expected ErrAccountLimitReached, got %v", err)
	}

	// Re-assigning an existing member must remain idempotent and must NOT
	// trip the cap (the row count does not grow).
	if err := mSvc.Assign(101, tenant.ID, true); err != nil {
		t.Fatalf("re-assign existing member should succeed, got %v", err)
	}
}

func TestAssign_UnlimitedWhenAccountLimitZero(t *testing.T) {
	setupTestDB(t)
	tSvc := tenantService{}
	mSvc := membershipService{}

	tenant, err := tSvc.Create(request.CreateTenantReq{
		Code:         "open",
		Name:         "Open Tenant",
		AccountLimit: 0,
	})
	if err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	for i := uint(1); i <= 5; i++ {
		if err := mSvc.Assign(200+i, tenant.ID, false); err != nil {
			t.Fatalf("assign user %d: %v", 200+i, err)
		}
	}

	members, err := mSvc.MembersOfTenant(tenant.ID)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if len(members) != 5 {
		t.Fatalf("expected 5 members under unlimited cap, got %d", len(members))
	}
}
