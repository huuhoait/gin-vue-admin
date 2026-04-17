package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// SysPolicyChangeLog is an append-only record of every Casbin policy
// mutation. Rows are never updated or deleted by the application; operators
// should revoke write access at the DB level for defense in depth so that a
// compromised admin cannot cover their tracks.
//
// PrevHash + Hash form a SHA-256 hash chain: each row's Hash covers all
// mutable fields of the row plus PrevHash of the preceding row. Tampering
// with any past row breaks the chain and is detectable by re-walking it.
type SysPolicyChangeLog struct {
	global.GVA_MODEL
	Actor       uint   `gorm:"index;comment:acting authority id"`
	ActorUserID uint   `gorm:"index;comment:acting user id"`
	Action      string `gorm:"size:32;index;comment:update|update_api|set_api_authorities|add|remove"`
	AuthorityID string `gorm:"size:64;index;comment:target role id if any"`
	Before      string `gorm:"type:text;comment:json snapshot before change"`
	After       string `gorm:"type:text;comment:json snapshot after change"`
	IP          string `gorm:"size:64;comment:source ip"`
	RequestID   string `gorm:"size:64;index;comment:correlation id"`
	Note        string `gorm:"size:255;comment:free-form description"`
	// Hash chain fields — do NOT update after insert.
	PrevHash string `gorm:"size:64;comment:hash of the previous row (hex SHA-256)"`
	Hash     string `gorm:"size:64;uniqueIndex;comment:SHA-256 of this row's content + PrevHash"`
}

func (SysPolicyChangeLog) TableName() string { return "sys_policy_change_logs" }
