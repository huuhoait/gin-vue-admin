package model

import (
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// OAuth2AuthCode is a single-use authorization code issued at /oauth2/authorize
// and exchanged for tokens at /oauth2/token. Default lifetime is 5 minutes.
type OAuth2AuthCode struct {
	global.GVA_MODEL
	Code        string     `json:"code" gorm:"size:128;uniqueIndex;not null"`
	ClientID    string     `json:"clientID" gorm:"size:64;index;not null"`
	UserID      uint       `json:"userID" gorm:"not null"`
	RedirectURI string     `json:"redirectURI" gorm:"type:text;not null"`
	Scope       string     `json:"scope" gorm:"type:text"`
	ExpiresAt   time.Time  `json:"expiresAt" gorm:"index"`
	UsedAt      *time.Time `json:"usedAt"`
}

func (OAuth2AuthCode) TableName() string { return "gva_oauth2_auth_codes" }
