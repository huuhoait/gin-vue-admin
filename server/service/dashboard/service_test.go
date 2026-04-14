package dashboard

import (
	"context"
	"testing"
)

func TestService_GetOverview_NilDBs(t *testing.T) {
	// With nil DBs, repository returns zeroes — no error.
	repo := NewRepository(nil, nil)
	cache := NewCache(nil)
	svc := NewService(repo, cache)

	m, err := svc.GetOverview(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.AgentsActive != 0 || m.OrdersToday != 0 || m.RevenueToday != 0 {
		t.Errorf("expected zeroes for nil DBs, got %+v", m)
	}
}

func TestRepository_NilDBs(t *testing.T) {
	repo := NewRepository(nil, nil)

	agents, err := repo.CountActiveAgents()
	if err != nil || agents != 0 {
		t.Errorf("CountActiveAgents: got %d, %v", agents, err)
	}

	orders, err := repo.CountOrdersToday()
	if err != nil || orders != 0 {
		t.Errorf("CountOrdersToday: got %d, %v", orders, err)
	}

	rev, err := repo.SumRevenueToday()
	if err != nil || rev != 0 {
		t.Errorf("SumRevenueToday: got %f, %v", rev, err)
	}
}
