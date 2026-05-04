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
	// PackageCode is an optional reference to gva_tenant_packages.code. The
	// link is intentionally a loose foreign key (no DB constraint) so that
	// deleting a package does not cascade and a tenant can outlive its
	// originally-assigned package while admins reconcile manually.
	PackageCode string `json:"packageCode" gorm:"size:64;index;comment:tenant package code"`
	Enabled     bool   `json:"enabled" gorm:"default:true;index"`
}

func (Tenant) TableName() string { return "gva_tenants" }
