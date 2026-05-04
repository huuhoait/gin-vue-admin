package model

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// TenantModel is the opt-in mixin for tenant-scoped models. Embed this
// instead of global.GVA_MODEL when a model's rows belong to exactly one
// tenant. The TenantID column is indexed for fast scoped queries.
//
// Usage:
//
//	type Order struct {
//	    tenantmodel.TenantModel
//	    Total int64
//	}
//
// Combined with service.WithTenantScope(ctx) in a query, rows are
// automatically filtered by the caller's tenant.
type TenantModel struct {
	global.GVA_MODEL
	TenantID uint `json:"tenantID" gorm:"index;not null;comment:owning tenant"`
}
