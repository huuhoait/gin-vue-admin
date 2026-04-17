package example

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

type ExaFileUploadAndDownload struct {
	global.GVA_MODEL
	Name    string `json:"name" form:"name" gorm:"column:name;comment:file name"`                                // file name
	ClassId int    `json:"classId" form:"classId" gorm:"default:0;type:int;column:class_id;comment:category ID;"` // category ID
	Url     string `json:"url" form:"url" gorm:"column:url;comment:file URL"`                                  // file URL
	Tag     string `json:"tag" form:"tag" gorm:"column:tag;comment:file tag"`                                  // file tag
	Key     string `json:"key" form:"key" gorm:"column:key;comment:serial number"`                                    // serial number
}

func (ExaFileUploadAndDownload) TableName() string {
	return "exa_file_upload_and_downloads"
}
