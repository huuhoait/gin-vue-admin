package model

import (
	"github.com/huuhoait/gin-vue-admin/server/global"

	"gorm.io/datatypes"
)

// TenantPackage bundles a curated set of SysBaseMenu IDs and SysApi IDs that
// can be referenced by a Tenant via Tenant.PackageCode. The package itself
// is a system-wide entity (NOT tenant-scoped) — it does not embed
// model.TenantModel and the tenant scoping callback ignores it.
//
// MenuIDs / ApiIDs are stored as JSON arrays of uint with a loose foreign key
// (no DB-level constraint). This keeps the package decoupled from gva_base_menus
// and sys_apis: deleting a menu/api will not cascade, but service-layer reads
// must tolerate stale IDs.
type TenantPackage struct {
	global.GVA_MODEL
	Code        string         `json:"code" gorm:"size:64;uniqueIndex;not null;comment:short stable identifier"`
	Name        string         `json:"name" gorm:"size:128;not null"`
	Description string         `json:"description" gorm:"type:text"`
	MenuIDs     datatypes.JSON `json:"menuIDs" gorm:"type:json;comment:array of SysBaseMenu IDs" swaggertype:"array,integer"`
	ApiIDs      datatypes.JSON `json:"apiIDs" gorm:"type:json;comment:array of SysApi IDs" swaggertype:"array,integer"`
	Enabled     bool           `json:"enabled" gorm:"default:true;index"`
}

func (TenantPackage) TableName() string { return "gva_tenant_packages" }
