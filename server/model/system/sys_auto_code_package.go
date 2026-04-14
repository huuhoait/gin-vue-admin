package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysAutoCodePackage struct {
	global.GVA_MODEL
	Desc        string `json:"desc" gorm:"comment:description"`
	Label       string `json:"label" gorm:"comment:display name"`
	Template    string `json:"template"  gorm:"comment:Template"`
	PackageName string `json:"packageName" gorm:"comment:package name"`
	Module      string `json:"-" example:"Module"`
}

func (s *SysAutoCodePackage) TableName() string {
	return "sys_auto_code_packages"
}
