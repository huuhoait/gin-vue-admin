package model

import (
	"github.com/huuhoait/gin-vue-admin/server/global"

	"gorm.io/datatypes"
)

// OAuth2Client is a registered third-party application that can request
// tokens. ClientSecret is never stored in plaintext — only the bcrypt hash
// is persisted; the plaintext is shown to the admin once at create/regen
// time and must be saved by the integrator.
type OAuth2Client struct {
	global.GVA_MODEL
	ClientID         string         `json:"clientID" gorm:"size:64;uniqueIndex;not null;comment:public client identifier"`
	ClientSecretHash string         `json:"-" gorm:"type:text;not null;comment:bcrypt(client_secret)"`
	Name             string         `json:"name" gorm:"size:128;not null"`
	Description      string         `json:"description" gorm:"type:text"`
	RedirectURIs     datatypes.JSON `json:"redirectUris" gorm:"type:json;comment:array of allowed redirect URIs" swaggertype:"array,string"`
	GrantTypes       datatypes.JSON `json:"grantTypes" gorm:"type:json;comment:authorization_code|refresh_token|client_credentials" swaggertype:"array,string"`
	Scopes           datatypes.JSON `json:"scopes" gorm:"type:json;comment:permitted scopes for this client" swaggertype:"array,string"`
	Enabled          bool           `json:"enabled" gorm:"default:true"`
	CreatedBy        uint           `json:"createdBy" gorm:"comment:admin user who created this client"`
}

func (OAuth2Client) TableName() string { return "gva_oauth2_clients" }
