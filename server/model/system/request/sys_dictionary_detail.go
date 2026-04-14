package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
	ParentID *uint `json:"parentID" form:"parentID"` // parent dictionary detail ID, Used forQuerySpecifyParentDownofSubItem
	Level    *int  `json:"level" form:"level"`       // hierarchy depth, Used forQuerySpecifyLayerLevelofData
}

// CreateSysDictionaryDetailRequest create dictionary detailRequest
type CreateSysDictionaryDetailRequest struct {
	Label           string `json:"label" form:"label" binding:"required"`                     // display value
	Value           string `json:"value" form:"value" binding:"required"`                     // dictionary value
	Extend          string `json:"extend" form:"extend"`                                      // extension value
	Status          *bool  `json:"status" form:"status"`                                      // enabled status
	Sort            int    `json:"sort" form:"sort"`                                          // sort order
	SysDictionaryID int    `json:"sysDictionaryID" form:"sysDictionaryID" binding:"required"` // association tag
	ParentID        *uint  `json:"parentID" form:"parentID"`                                  // parent dictionary detail ID
}

// UpdateSysDictionaryDetailRequest update dictionary detailRequest
type UpdateSysDictionaryDetailRequest struct {
	ID              uint   `json:"ID" form:"ID" binding:"required"`                           // primary key ID
	Label           string `json:"label" form:"label" binding:"required"`                     // display value
	Value           string `json:"value" form:"value" binding:"required"`                     // dictionary value
	Extend          string `json:"extend" form:"extend"`                                      // extension value
	Status          *bool  `json:"status" form:"status"`                                      // enabled status
	Sort            int    `json:"sort" form:"sort"`                                          // sort order
	SysDictionaryID int    `json:"sysDictionaryID" form:"sysDictionaryID" binding:"required"` // association tag
	ParentID        *uint  `json:"parentID" form:"parentID"`                                  // parent dictionary detail ID
}

// GetDictionaryDetailsByParentRequest according to parentIDget dictionaryDetailsRequest
type GetDictionaryDetailsByParentRequest struct {
	SysDictionaryID int   `json:"sysDictionaryID" form:"sysDictionaryID" binding:"required"` // DictionaryID
	ParentID        *uint `json:"parentID" form:"parentID"`                                  // parent dictionary detail ID, When emptygetTopLevel
	IncludeChildren bool  `json:"includeChildren" form:"includeChildren"`                    // YesNoPackageIncludeSubLevelData
}
