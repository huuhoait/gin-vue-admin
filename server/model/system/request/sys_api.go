package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// apiPaginationConditionQueryAndsortStructureBody
type SearchApiParams struct {
	system.SysApi
	request.PageInfo
	OrderKey string `json:"orderKey"` // sort
	Desc     bool   `json:"desc"`     // sortMethod:Ascendingfalse(default)|Descendingtrue
}

// SetApiAuthorities ApprovedAPIpathAndmethodfull overwriteAssociationRoleList
type SetApiAuthorities struct {
	Path         string `json:"path" form:"path"`                     // APIpath
	Method       string `json:"method" form:"method"`                 // request method
	AuthorityIds []uint `json:"authorityIds" form:"authorityIds"`     // role IDList
}
