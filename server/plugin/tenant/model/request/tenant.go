package request

import "github.com/huuhoait/gin-vue-admin/server/model/common/request"

type CreateTenantReq struct {
	Code        string `json:"code" binding:"required,min=2,max=64"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateTenantReq struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     *bool  `json:"enabled"`
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
