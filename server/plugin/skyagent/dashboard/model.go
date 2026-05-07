package dashboard

// Minimal read-only GORM models for aggregate queries.
// These do NOT represent the full schema — only fields needed for COUNT/SUM.

// AgentRow maps to the agents table for aggregate queries.
type AgentRow struct {
	Status string `gorm:"column:status"`
}

func (AgentRow) TableName() string { return "agents" }

// OrderRow maps to the orders table for aggregate queries.
type OrderRow struct {
	TotalAmount float64 `gorm:"column:total_amount"`
}

func (OrderRow) TableName() string { return "orders" }

// OverviewMetrics is the response payload for the dashboard overview endpoint.
type OverviewMetrics struct {
	AgentsActive int64   `json:"agents_active"`
	OrdersToday  int64   `json:"orders_today"`
	RevenueToday float64 `json:"revenue_today"`
}
