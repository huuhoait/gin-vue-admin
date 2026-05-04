package model

import (
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// Tenant represents a logical isolation boundary (organization, workspace,
// customer account). Models opt into tenant scoping by embedding TenantModel.
//
// Lifecycle attributes (mirrors ruoyi-vue-pro tenant model):
//   - Enabled: hard switch; disabled tenants reject login/membership use.
//   - ExpireAt: optional expiration timestamp; nil means "never expires".
//   - AccountLimit: cap on the number of users that can be assigned to this
//     tenant. 0 means "unlimited". Enforced by the membership service when
//     assigning new members; existing members beyond the cap are not pruned.
//
// Reads (FindByID, FindByCode) intentionally do NOT short-circuit on
// disabled/expired — callers (login flow, tenant middleware) decide policy
// using IsActive(t) so admin tooling can still load and edit a disabled
// tenant.
type Tenant struct {
	global.GVA_MODEL
	Code         string     `json:"code" gorm:"size:64;uniqueIndex;not null;comment:short stable identifier"`
	Name         string     `json:"name" gorm:"size:128;not null"`
	Description  string     `json:"description" gorm:"type:text"`
	ContactName  string     `json:"contactName" gorm:"size:128;comment:contact person name"`
	ContactPhone string     `json:"contactPhone" gorm:"size:64;comment:contact phone number"`
	Domain       string     `json:"domain" gorm:"size:255;index;comment:subdomain code or vanity domain"`
	ExpireAt     *time.Time `json:"expireAt" gorm:"comment:expiration time; null = never"`
	AccountLimit int        `json:"accountLimit" gorm:"default:0;comment:max member count; 0 = unlimited"`
	Enabled      bool       `json:"enabled" gorm:"default:true;index"`
}

func (Tenant) TableName() string { return "gva_tenants" }
