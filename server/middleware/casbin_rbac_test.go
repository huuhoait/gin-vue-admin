package middleware

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func TestSkyAgentRBAC_KycReviewerDeniedCommission(t *testing.T) {
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

	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("new enforcer: %v", err)
	}

	// Seeded baseline from source/system/casbin.go
	_, _ = e.AddPolicy("9100", "/kyc/ping", "GET")
	_, _ = e.AddPolicy("9200", "/commission/ping", "GET")

	ok, err := e.Enforce("9100", "/kyc/ping", "GET")
	if err != nil {
		t.Fatalf("enforce kyc: %v", err)
	}
	if !ok {
		t.Fatalf("expected allow for KYC reviewer on /kyc/ping")
	}

	ok, err = e.Enforce("9100", "/commission/ping", "GET")
	if err != nil {
		t.Fatalf("enforce commission: %v", err)
	}
	if ok {
		t.Fatalf("expected deny for KYC reviewer on /commission/ping")
	}
}

