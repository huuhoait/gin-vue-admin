package dashboard

import (
	"time"

	"gorm.io/gorm"
)

// Repository provides aggregate read-only queries against Core and Order databases.
type Repository struct {
	coreDB  *gorm.DB
	orderDB *gorm.DB
}

// NewRepository creates a dashboard Repository. Either DB may be nil if
// the connection is unavailable — the corresponding metric will return 0.
func NewRepository(coreDB, orderDB *gorm.DB) *Repository {
	return &Repository{coreDB: coreDB, orderDB: orderDB}
}

// CountActiveAgents returns the number of agents with status "active".
func (r *Repository) CountActiveAgents() (int64, error) {
	if r.coreDB == nil {
		return 0, nil
	}
	var count int64
	err := r.coreDB.Model(&AgentRow{}).Where("status = ?", "active").Count(&count).Error
	return count, err
}

// CountOrdersToday returns the number of orders created since today 00:00 UTC.
func (r *Repository) CountOrdersToday() (int64, error) {
	if r.orderDB == nil {
		return 0, nil
	}
	todayUTC := time.Now().UTC().Truncate(24 * time.Hour)
	var count int64
	err := r.orderDB.Model(&OrderRow{}).Where("created_at >= ?", todayUTC).Count(&count).Error
	return count, err
}

// SumRevenueToday returns the total revenue (sum of total_amount) for orders created today.
func (r *Repository) SumRevenueToday() (float64, error) {
	if r.orderDB == nil {
		return 0, nil
	}
	todayUTC := time.Now().UTC().Truncate(24 * time.Hour)
	var sum float64
	err := r.orderDB.Model(&OrderRow{}).
		Where("created_at >= ?", todayUTC).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&sum).Error
	return sum, err
}

// FetchOverview aggregates all dashboard metrics in a single call.
func (r *Repository) FetchOverview() (*OverviewMetrics, error) {
	agents, err := r.CountActiveAgents()
	if err != nil {
		return nil, err
	}
	orders, err := r.CountOrdersToday()
	if err != nil {
		return nil, err
	}
	revenue, err := r.SumRevenueToday()
	if err != nil {
		return nil, err
	}
	return &OverviewMetrics{
		AgentsActive: agents,
		OrdersToday:  orders,
		RevenueToday: revenue,
	}, nil
}
