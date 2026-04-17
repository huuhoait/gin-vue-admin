// auto-generate templateSysParams
package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// Parameter StructureBody  SysParams
type SysParams struct {
	global.GVA_MODEL
	Name  string `json:"name" form:"name" gorm:"column:name;comment:Parametername;" binding:"required"`   //Parametername
	Key   string `json:"key" form:"key" gorm:"column:key;comment:parameter key;" binding:"required"`       //parameter key
	Value string `json:"value" form:"value" gorm:"column:value;comment:parameter value;" binding:"required"` //parameter value
	Desc  string `json:"desc" form:"desc" gorm:"column:desc;comment:parameter description;"`                      //parameter description
}

// TableName Parameter SysParamsCustomtable name sys_params
func (SysParams) TableName() string {
	return "sys_params"
}
