package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
)

// setupTestDB swaps global.GVA_DB for an in-memory SQLite database with the
// gva_user_tenants schema migrated. Returns a cleanup that restores the
// previous DB pointer so suites running in parallel don't bleed state.
func setupTestDB(t *testing.T) func() {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.UserTenant{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	prev := global.GVA_DB
	global.GVA_DB = db
	return func() { global.GVA_DB = prev }
}

// withClaims registers a middleware that mimics middleware.JWTAuth by
// stashing the supplied custom claims onto the request context. The tenant
// middleware reads claims via utils.GetUserID / utils.GetTenantID which
// inspect c.Get("claims") first, so this is enough to drive the unit
// without instantiating a JWT.
func withClaims(claims *systemReq.CustomClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	}
}

// runRequest fires a single request through a router that has both the
// fake JWT middleware and the real TenantContext middleware mounted, then
// returns the resolved tenant id (0 when unscoped) and the response.
func runRequest(t *testing.T, claims *systemReq.CustomClaims, header string) (resolvedTID uint, scoped bool, status int, body string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(withClaims(claims))
	r.Use(TenantContext())
	r.GET("/probe", func(c *gin.Context) {
		if id, ok := service.FromContext(c.Request.Context()); ok {
			resolvedTID = id
			scoped = true
		}
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	if header != "" {
		req.Header.Set("X-Tenant-ID", header)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return resolvedTID, scoped, w.Code, w.Body.String()
}

func makeClaims(userID, tenantID uint) *systemReq.CustomClaims {
	return &systemReq.CustomClaims{
		BaseClaims: systemReq.BaseClaims{
			UUID:     uuid.New(),
			ID:       userID,
			Username: "alice",
			TenantID: tenantID,
		},
	}
}

// TestTenantContext_FromJWT_NoDB asserts that a JWT carrying tenant_id=5
// makes it onto the request context with no membership lookup. We assert
// the no-DB property by leaving global.GVA_DB nil — any DB call would
// nil-panic, which the test would catch.
func TestTenantContext_FromJWT_NoDB(t *testing.T) {
	prev := global.GVA_DB
	global.GVA_DB = nil
	t.Cleanup(func() { global.GVA_DB = prev })

	tid, scoped, status, _ := runRequest(t, makeClaims(42, 5), "")
	if status != http.StatusOK {
		t.Fatalf("status=%d, want 200", status)
	}
	if !scoped || tid != 5 {
		t.Fatalf("resolved tenant=%d scoped=%v, want 5/true", tid, scoped)
	}
}

// TestTenantContext_HeaderOverride_ZeroClaim covers the common
// super-admin case: JWT carries tenant_id=0 (system tenant), header asks
// for tenant 5, membership exists.
func TestTenantContext_HeaderOverride_ZeroClaim(t *testing.T) {
	cleanup := setupTestDB(t)
	t.Cleanup(cleanup)

	if err := global.GVA_DB.Create(&model.UserTenant{UserID: 42, TenantID: 5, IsPrimary: true}).Error; err != nil {
		t.Fatalf("seed membership: %v", err)
	}

	tid, scoped, status, _ := runRequest(t, makeClaims(42, 0), "5")
	if status != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", status, "")
	}
	if !scoped || tid != 5 {
		t.Fatalf("resolved tenant=%d scoped=%v, want 5/true", tid, scoped)
	}
}

// TestTenantContext_HeaderOverride_DenyOnMissingMembership asserts that
// when a non-zero JWT tenant disagrees with a header the middleware
// performs the membership check and 403s when the user has no row.
func TestTenantContext_HeaderOverride_DenyOnMissingMembership(t *testing.T) {
	cleanup := setupTestDB(t)
	t.Cleanup(cleanup)

	// Membership exists for tenant 5 only.
	if err := global.GVA_DB.Create(&model.UserTenant{UserID: 42, TenantID: 5, IsPrimary: true}).Error; err != nil {
		t.Fatalf("seed membership: %v", err)
	}

	_, scoped, status, body := runRequest(t, makeClaims(42, 5), "7")
	if scoped {
		t.Fatalf("expected the request to be aborted, but tenant scope was applied")
	}
	if status != http.StatusOK { // FailWithCode returns 200 with code=7 envelope
		t.Fatalf("status=%d, want 200", status)
	}
	if body == "" {
		t.Fatalf("expected error envelope body")
	}
}

// TestTenantContext_SystemTenant_NoHeader asserts JWT tenant_id=0 +
// no header leaves the context unscoped (FromContext returns ok=false).
// This is the super-admin "see everything" path. No DB access required.
func TestTenantContext_SystemTenant_NoHeader(t *testing.T) {
	prev := global.GVA_DB
	global.GVA_DB = nil
	t.Cleanup(func() { global.GVA_DB = prev })

	tid, scoped, status, _ := runRequest(t, makeClaims(42, 0), "")
	if status != http.StatusOK {
		t.Fatalf("status=%d, want 200", status)
	}
	if scoped || tid != 0 {
		t.Fatalf("expected unscoped context, got tenant=%d scoped=%v", tid, scoped)
	}
}
