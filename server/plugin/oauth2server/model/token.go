package model

import (
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// OAuth2Token holds an access/refresh pair. Opaque random strings, not JWTs —
// keeps revocation cheap (one row update) and avoids leaking signing keys to
// clients. UserID is 0 for client_credentials grants.
type OAuth2Token struct {
	global.GVA_MODEL
	AccessToken      string     `json:"-" gorm:"size:128;uniqueIndex;not null"`
	RefreshToken     string     `json:"-" gorm:"size:128;uniqueIndex"`
	ClientID         string     `json:"clientID" gorm:"size:64;index;not null"`
	UserID           uint       `json:"userID" gorm:"index;comment:0 for client_credentials"`
	Scope            string     `json:"scope" gorm:"type:text"`
	GrantType        string     `json:"grantType" gorm:"size:32"`
	AccessExpiresAt  time.Time  `json:"accessExpiresAt"`
	RefreshExpiresAt *time.Time `json:"refreshExpiresAt"`
	RevokedAt        *time.Time `json:"revokedAt"`
}

func (OAuth2Token) TableName() string { return "gva_oauth2_tokens" }
