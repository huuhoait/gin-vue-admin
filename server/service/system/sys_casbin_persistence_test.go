package system

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestCasbinPoliciesPersistToCasbinRule(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&gormadapter.CasbinRule{}); err != nil {
		t.Fatalf("migrate casbin_rule: %v", err)
	}

	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		t.Fatalf("build model: %v", err)
	}

	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		t.Fatalf("new gorm adapter: %v", err)
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		t.Fatalf("new enforcer: %v", err)
	}

	added, err := e.AddPolicy("9100", "/kyc/ping", "GET")
	if err != nil {
		t.Fatalf("add policy: %v", err)
	}
	if !added {
		t.Fatalf("expected policy to be added")
	}

	var count int64
	if err := db.Model(&gormadapter.CasbinRule{}).Count(&count).Error; err != nil {
		t.Fatalf("count casbin_rule: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected casbin_rule rows >= 1, got %d", count)
	}
}

