package model

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// Tenant represents a logical isolation boundary (organization, workspace,
// customer account). Models opt into tenant scoping by embedding TenantModel.
type Tenant struct {
	global.GVA_MODEL
	Code        string `json:"code" gorm:"size:64;uniqueIndex;not null;comment:short stable identifier"`
	Name        string `json:"name" gorm:"size:128;not null"`
	Description string `json:"description" gorm:"type:text"`
	Enabled     bool   `json:"enabled" gorm:"default:true;index"`
}

func (Tenant) TableName() string { return "gva_tenants" }
