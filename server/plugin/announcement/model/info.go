package model

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// Info Announcement StructureBody
type Info struct {
	global.GVA_MODEL
	Title       string         `json:"title" form:"title" gorm:"column:title;comment:announcement title;"`                                              //Title
	Content     string         `json:"content" form:"content" gorm:"column:content;comment:announcement content;type:text;"`                              //content
	UserID      *int           `json:"userID" form:"userID" gorm:"column:user_id;comment:publisher;"`                                           //DoPerson
	Attachments datatypes.JSON `json:"attachments" form:"attachments" gorm:"column:attachments;comment:attachments;" swaggertype:"array,object"` //Attachment
}

// TableName Announcement InfoCustomtable name gva_announcements_info
func (Info) TableName() string {
	return "gva_announcements_info"
}
