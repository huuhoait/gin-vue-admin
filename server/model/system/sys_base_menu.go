package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

type SysBaseMenu struct {
	global.GVA_MODEL
	MenuLevel     uint                   `json:"-"`
	ParentId      uint                   `json:"parentId" gorm:"comment:Parentmenu ID"`          // Parentmenu ID
	Path          string                 `json:"path" gorm:"comment:route path"`              // route path
	Name          string                 `json:"name" gorm:"comment:route name"`              // route name
	Hidden        bool                   `json:"hidden" gorm:"comment:YesNoin listhidden"`      // YesNoin listhidden
	Component     string                 `json:"component" gorm:"comment:corresponding frontend filepath"` // corresponding frontend filepath
	Sort          int                    `json:"sort" gorm:"comment:sort order"`              // sort order
	Meta          `json:"meta" gorm:"embedded"`                                             // additional properties
	SysAuthoritys []SysAuthority         `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
	Children      []SysBaseMenu          `json:"children" gorm:"-"`
	Parameters    []SysBaseMenuParameter `json:"parameters"`
	MenuBtn       []SysBaseMenuBtn       `json:"menuBtn"`
}

type Meta struct {
	ActiveName     string `json:"activeName" gorm:"comment:HighBrightMenu"`
	KeepAlive      bool   `json:"keepAlive" gorm:"comment:YesNoCache"`                 // YesNoCache
	DefaultMenu    bool   `json:"defaultMenu" gorm:"comment:YesNoYesbase route (dev)"` // YesNoYesbase route (dev)
	Title          string `json:"title" gorm:"comment:menu title"`                       // menu title
	Icon           string `json:"icon" gorm:"comment:MenuIcon"`                      // MenuIcon
	CloseTab       bool   `json:"closeTab" gorm:"comment:AutomaticDisabletab"`               // AutomaticDisabletab
	TransitionType string `json:"transitionType" gorm:"comment:RouteSwitchMovePaint"`        // RouteSwitchMovePaint
}

type SysBaseMenuParameter struct {
	global.GVA_MODEL
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:address bar query param isparamsStillYesquery"` // address bar query param isparamsStillYesquery
	Key           string `json:"key" gorm:"comment:query string in address barkey"`              // query string in address barkey
	Value         string `json:"value" gorm:"comment:value of address bar query param"`             // value of address bar query param
}

func (SysBaseMenu) TableName() string {
	return "sys_base_menus"
}
