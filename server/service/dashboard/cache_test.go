package dashboard

import (
	"context"
	"testing"
)

func TestCache_NilRedis(t *testing.T) {
	c := NewCache(nil)

	// Get on nil redis should return nil (cache miss).
	if m := c.Get(context.Background()); m != nil {
		t.Errorf("expected nil on nil redis, got %+v", m)
	}

	// Set on nil redis should not panic.
	c.Set(context.Background(), &OverviewMetrics{AgentsActive: 1})
}

func TestOverviewMetrics_JSON(t *testing.T) {
	m := OverviewMetrics{
		AgentsActive: 42,
		OrdersToday:  100,
		RevenueToday: 5000000,
	}
	if m.AgentsActive != 42 || m.OrdersToday != 100 || m.RevenueToday != 5000000 {
		t.Errorf("unexpected metrics: %+v", m)
	}
}
