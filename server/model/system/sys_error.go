// auto-generate templateSysError
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Error Log StructureBody  SysError
type SysError struct {
	global.GVA_MODEL
	Form     *string `json:"form" form:"form" gorm:"comment:ErrorSource;column:form;type:text;" binding:"required"` //ErrorSource
	Info     *string `json:"info" form:"info" gorm:"comment:Errorcontent;column:info;type:text;"`                    //Errorcontent
	Level    string  `json:"level" form:"level" gorm:"comment:LogLevel;column:level;"`
	Solution *string `json:"solution" form:"solution" gorm:"comment:ResolvePlan;column:solution;type:text"`               //ResolvePlan
	Status   string  `json:"status" form:"status" gorm:"comment:Handlestatus;column:status;type:varchar(20);default:NotHandle;"` //Handlestatus:NotHandle/HandleIn/HandleComplete
}

// TableName Error Log SysErrorCustomtable name sys_error
func (SysError) TableName() string {
	return "sys_error"
}
