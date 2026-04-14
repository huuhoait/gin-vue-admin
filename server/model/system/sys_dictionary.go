// auto-generate templateSysDictionary
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// If it contains time.Time please import the time package yourself
type SysDictionary struct {
	global.GVA_MODEL
	Name                 string                `json:"name" form:"name" gorm:"column:name;comment:dictionary name"`              // dictionary name
	Type                 string                `json:"type" form:"type" gorm:"column:type;comment:dictionary name(EN)"`              // dictionary name(EN)
	Status               *bool                 `json:"status" form:"status" gorm:"column:status;comment:status"`            // status
	Desc                 string                `json:"desc" form:"desc" gorm:"column:desc;comment:description"`                  // description
	ParentID             *uint                 `json:"parentID" form:"parentID" gorm:"column:parent_id;comment:ParentDictionaryID"` // ParentDictionaryID
	Children             []SysDictionary       `json:"children" gorm:"foreignKey:ParentID"`                             // SubDictionary
	SysDictionaryDetails []SysDictionaryDetail `json:"sysDictionaryDetails" form:"sysDictionaryDetails"`
}

func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}
