// auto-generate templateSysVersion
package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// Version Management StructureBody  SysVersion
type SysVersion struct {
	global.GVA_MODEL
	VersionName *string `json:"versionName" form:"versionName" gorm:"comment:version name;column:version_name;size:255;" binding:"required"`                           //version name
	VersionCode *string `json:"versionCode" form:"versionCode" gorm:"comment:version number;column:version_code;size:100;" binding:"required"`                            //version number
	Description *string `json:"description" form:"description" gorm:"comment:version description;column:description;size:500;"`                                               //version description
	VersionData *string `json:"versionData" form:"versionData" gorm:"comment:VersionDataJSON;column:version_data;type:text;"` //VersionData
}

// TableName Version Management SysVersionCustomtable name sys_versions
func (SysVersion) TableName() string {
	return "sys_versions"
}
