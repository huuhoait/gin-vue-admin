package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

type JwtBlacklist struct {
	global.GVA_MODEL
	// Jwt is stored as text (tokens are > 255 chars) but we still index a
	// prefix so audit lookups by full token stay fast; MySQL/PG will pick up
	// the hash-based unique path automatically on equality lookups.
	Jwt string `gorm:"type:text;comment:jwt;index:idx_jwt_blacklist_jwt,length:128"`
}
