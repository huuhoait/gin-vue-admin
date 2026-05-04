package global

import (
	"time"

	"gorm.io/gorm"
)

// GVA_MODEL is the base composite embedded by every persistent model in the
// project. It provides the conventional ID/timestamps/soft-delete plus a
// basic audit trail (CreatedBy/UpdatedBy/DeletedBy) populated automatically
// by the audit GORM callback when the request context carries a user id.
//
// The audit fields default to 0, which is the convention for "system origin"
// — used for rows created by source/system/* seeds (no request context),
// background cron jobs, and migrations. Existing tables migrate cleanly:
// AutoMigrate adds the columns with default 0, no row rewrites required.
type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey" json:"ID"`
	CreatedAt time.Time      // created at — JSON key stays "CreatedAt" (Go field name) for FE compat
	UpdatedAt time.Time      // updated at — same
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy uint           `gorm:"default:0;comment:user id who created this row" json:"createdBy"`
	UpdatedBy uint           `gorm:"default:0;comment:user id who last updated this row" json:"updatedBy"`
	DeletedBy uint           `gorm:"default:0;comment:user id who soft-deleted this row" json:"deletedBy"`
}
