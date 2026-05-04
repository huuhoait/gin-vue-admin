package model

import "time"

// UserTenant links a SysUser (referenced loosely by ID — no FK to avoid
// coupling the plugin to the system user table) to a Tenant. Composite
// primary key (UserID, TenantID).
//
// IsPrimary marks the tenant the user is placed into when no X-Tenant-ID
// header is sent. Exactly one row per user should have IsPrimary=true; this
// is a soft invariant enforced by the service layer.
type UserTenant struct {
	UserID    uint      `json:"userID" gorm:"primaryKey"`
	TenantID  uint      `json:"tenantID" gorm:"primaryKey;index"`
	IsPrimary bool      `json:"isPrimary" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserTenant) TableName() string { return "gva_user_tenants" }
