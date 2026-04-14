package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysApi struct {
	global.GVA_MODEL
	Path        string `json:"path" gorm:"comment:API path"`             // API path
	Description string `json:"description" gorm:"comment:API description"`    // API description
	ApiGroup    string `json:"apiGroup" gorm:"comment:API group"`          // API group
	Method      string `json:"method" gorm:"default:POST;comment:method"` // method:CreatePOST(default)|viewGET|updatePUT|deleteDELETE
}

func (SysApi) TableName() string {
	return "sys_apis"
}

type SysIgnoreApi struct {
	global.GVA_MODEL
	Path   string `json:"path" gorm:"comment:API path"`             // API path
	Method string `json:"method" gorm:"default:POST;comment:method"` // method:CreatePOST(default)|viewGET|updatePUT|deleteDELETE
	Flag   bool   `json:"flag" gorm:"-"`                         // YesNoIgnore
}

func (SysIgnoreApi) TableName() string {
	return "sys_ignore_apis"
}
