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
}

// UpdateTenantReq is the payload for PUT /tenant/update. All optional fields
// use pointer / sentinel semantics so that omitted values do not overwrite
// existing data:
//   - String fields: empty string means "do not change" (preserves prior
//     behaviour). Use UpdateTenantClearReq variant if explicit clearing is
//     needed in future.
//   - ExpireAt:     nil pointer means "do not change". To clear an
//     expiration explicitly, use ClearExpireAt=true.
//   - AccountLimit: pointer; nil means "do not change". Pass 0 explicitly
//     to mark as unlimited.
//   - Enabled:      pointer; nil means "do not change".
type UpdateTenantReq struct {
	ID             uint       `json:"id" binding:"required"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	ContactName    string     `json:"contactName"`
	ContactPhone   string     `json:"contactPhone"`
	Domain         string     `json:"domain"`
	ExpireAt       *time.Time `json:"expireAt"`
	ClearExpireAt  bool       `json:"clearExpireAt"`
	AccountLimit   *int       `json:"accountLimit" binding:"omitempty,min=0"`
	Enabled        *bool      `json:"enabled"`
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

type UnassignUserReq struct {
	UserID   uint `json:"userID" binding:"required"`
	TenantID uint `json:"tenantID" binding:"required"`
}
