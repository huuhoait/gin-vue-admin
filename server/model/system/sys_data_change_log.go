package system

import "github.com/huuhoait/gin-vue-admin/server/global"

// SysDataChangeLog records before/after snapshots for every admin-initiated
// mutation to user records, system configuration, and other sensitive entities.
// Rows are append-only; the application never updates or deletes them.
type SysDataChangeLog struct {
	global.GVA_MODEL
	// Who made the change
	OperatorID   uint   `gorm:"index;comment:user id of the operator"`
	OperatorName string `gorm:"size:64;comment:username of the operator"`
	// What was changed
	TargetType string `gorm:"size:64;index;comment:entity type e.g. SysUser / SystemConfig"`
	TargetID   string `gorm:"size:64;index;comment:primary key or identifier of the target"`
	Action     string `gorm:"size:32;index;comment:create|update|delete|reset_password|set_authority"`
	// Snapshots (JSON)
	Before string `gorm:"type:text;comment:json snapshot before change"`
	After  string `gorm:"type:text;comment:json snapshot after change"`
	// Request context
	IP        string `gorm:"size:64;comment:operator source IP"`
	RequestID string `gorm:"size:64;index;comment:X-Request-ID correlation id"`
	// Tamper-evident hash chain
	PrevHash string `gorm:"size:64;comment:hash of the previous row (hex SHA-256)"`
	Hash     string `gorm:"size:64;uniqueIndex;comment:SHA-256 of this row content + PrevHash"`
}

func (SysDataChangeLog) TableName() string { return "sys_data_change_logs" }
