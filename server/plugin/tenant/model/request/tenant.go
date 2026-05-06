package request

import (
	"time"

	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
)

// CreateTenantReq is the payload for POST /tenant/create.
//
// ExpireAt is a pointer so the FE can omit / null it for "never expires".
// AccountLimit defaults to 0 (unlimited) if omitted.
type CreateTenantReq struct {
	Code         string     `json:"code" binding:"required,min=2,max=64"`
	Name         string     `json:"name" binding:"required"`
	Description  string     `json:"description"`
	ContactName  string     `json:"contactName"`
	ContactPhone string     `json:"contactPhone"`
	Domain       string     `json:"domain"`
	ExpireAt     *time.Time `json:"expireAt"`
	AccountLimit int        `json:"accountLimit" binding:"omitempty,min=0"`
	PackageCode  string     `json:"packageCode"`
}

// UpdateTenantReq is the payload for PUT /tenant/update. Optional fields use
// pointer / sentinel semantics so omitted values don't overwrite existing
// data:
//   - String fields: empty string means "do not change".
//   - ExpireAt: nil pointer means "do not change". Use ClearExpireAt=true to
//     explicitly clear an expiration.
//   - AccountLimit: pointer; nil means "do not change". 0 = unlimited.
//   - Enabled: pointer; nil means "do not change".
type UpdateTenantReq struct {
	ID            uint       `json:"id" binding:"required"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	ContactName   string     `json:"contactName"`
	ContactPhone  string     `json:"contactPhone"`
	Domain        string     `json:"domain"`
	ExpireAt      *time.Time `json:"expireAt"`
	ClearExpireAt bool       `json:"clearExpireAt"`
	AccountLimit  *int       `json:"accountLimit" binding:"omitempty,min=0"`
	PackageCode   string     `json:"packageCode"`
	Enabled       *bool      `json:"enabled"`
}

type TenantListReq struct {
	request.PageInfo
	Enabled *bool `json:"enabled" form:"enabled"`
}

type IdReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type AssignUserReq struct {
	UserID    uint `json:"userID" binding:"required"`
	TenantID  uint `json:"tenantID" binding:"required"`
	IsPrimary bool `json:"isPrimary"`
}

// UnassignUserReq is bound from the query string (DELETE) — Gin's query
// binder uses the form: tag, so json: alone wouldn't match.
type UnassignUserReq struct {
	UserID   uint `json:"userID" form:"userID" binding:"required"`
	TenantID uint `json:"tenantID" form:"tenantID" binding:"required"`
}

// CreateUserAndAssignReq creates a fresh SysUser scoped to a tenant in one
// step. Username + Password are required; the rest are optional profile bits.
// The new user is given the default "Tenant" authority (id 9300) and added
// to the target tenant's membership in the same transaction so the tenant
// AccountLimit check is enforced atomically.
type CreateUserAndAssignReq struct {
	TenantID  uint   `json:"tenantID" binding:"required"`
	Username  string `json:"userName" binding:"required,min=3,max=64"`
	Password  string `json:"password" binding:"required,min=6,max=128"`
	NickName  string `json:"nickName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	IsPrimary bool   `json:"isPrimary"`
}
