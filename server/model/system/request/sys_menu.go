package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// AddMenuAuthorityInfo Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []system.SysBaseMenu `json:"menus"`
	AuthorityId uint                 `json:"authorityId"` // role ID
}

// SetMenuAuthorities Approvedmenu IDfull overwriteAssociationRoleList
type SetMenuAuthorities struct {
	MenuId       uint   `json:"menuId" form:"menuId"`             // menu ID
	AuthorityIds []uint `json:"authorityIds" form:"authorityIds"` // role IDList
}

func DefaultMenu() []system.SysBaseMenu {
	return []system.SysBaseMenu{{
		GVA_MODEL: global.GVA_MODEL{ID: 1},
		ParentId:  0,
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Meta: system.Meta{
			Title: "Dashboard",
			Icon:  "setting",
		},
	}}
}
