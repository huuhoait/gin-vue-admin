package request

import "github.com/huuhoait/gin-vue-admin/server/model/common/request"

// CreatePackageReq is the body shape for POST /tenantPackage/create.
// MenuIDs/ApiIDs are passed as plain []uint and serialized to JSON in the
// service layer.
type CreatePackageReq struct {
	Code        string `json:"code" binding:"required,min=2,max=64"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MenuIDs     []uint `json:"menuIDs"`
	ApiIDs      []uint `json:"apiIDs"`
}

// UpdatePackageReq mutates an existing package by ID. Code is immutable; if
// callers send a new code it is silently ignored. Slice fields are nullable —
// a nil slice means "leave unchanged"; an empty slice means "clear".
type UpdatePackageReq struct {
	ID          uint    `json:"id" binding:"required"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MenuIDs     *[]uint `json:"menuIDs"`
	ApiIDs      *[]uint `json:"apiIDs"`
	Enabled     *bool   `json:"enabled"`
}

// PackageListReq drives /tenantPackage/list. Reuses the common PageInfo for
// page/pageSize/keyword.
type PackageListReq struct {
	request.PageInfo
	Enabled *bool `json:"enabled" form:"enabled"`
}
