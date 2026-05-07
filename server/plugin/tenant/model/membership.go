package model

import "time"

// UserTenant links a SysUser (referenced loosely by ID — no FK to avoid
// coupling the plugin to the system user table) to a Tenant. Composite
// primary key (UserID, TenantID).
//
// IsPrimary marks the tenant the user is placed into when no X-Tenant-ID
// header is sent. Exactly one row per user should have IsPrimary=true; this
// is a soft invariant enforced by the service layer.
// Composite index idx_user_tenants_tenant_primary_created backs MembersOfTenant's
// "WHERE tenant_id = ? ORDER BY is_primary DESC, created_at ASC" so Postgres
// can serve it from an indexed scan instead of a sort. Column priorities match
// the ORDER BY tail; the ASC default on created_at is fine because the leading
// is_primary DESC determines sort direction usability.
type UserTenant struct {
	UserID    uint      `json:"userID" gorm:"primaryKey"`
	TenantID  uint      `json:"tenantID" gorm:"primaryKey;index;index:idx_user_tenants_tenant_primary_created,priority:1"`
	IsPrimary bool      `json:"isPrimary" gorm:"default:false;index:idx_user_tenants_tenant_primary_created,sort:desc,priority:2"`
	CreatedAt time.Time `json:"createdAt" gorm:"index:idx_user_tenants_tenant_primary_created,priority:3"`
}

func (UserTenant) TableName() string { return "gva_user_tenants" }
