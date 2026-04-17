// auto-generate templateSysDictionaryDetail
package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// If it contains time.Time please import the time package yourself
type SysDictionaryDetail struct {
	global.GVA_MODEL
	Label           string                `json:"label" form:"label" gorm:"column:label;comment:display value"`                                  // display value
	Value           string                `json:"value" form:"value" gorm:"column:value;comment:dictionary value"`                                  // dictionary value
	Extend          string                `json:"extend" form:"extend" gorm:"column:extend;comment:extension value"`                               // extension value
	Status          *bool                 `json:"status" form:"status" gorm:"column:status;comment:enabled status"`                              // enabled status
	Sort            int                   `json:"sort" form:"sort" gorm:"column:sort;comment:sort order"`                                    // sort order
	SysDictionaryID int                   `json:"sysDictionaryID" form:"sysDictionaryID" gorm:"column:sys_dictionary_id;comment:association tag"` // association tag
	ParentID        *uint                 `json:"parentID" form:"parentID" gorm:"column:parent_id;comment:parent dictionary detail ID"`                   // parent dictionary detail ID
	Children        []SysDictionaryDetail `json:"children" gorm:"foreignKey:ParentID"`                                                 // SubDictionaryDetails
	Level           int                   `json:"level" form:"level" gorm:"column:level;comment:hierarchy depth"`                                 // hierarchy depth, From0Start
	Path            string                `json:"path" form:"path" gorm:"column:path;comment:LayerLevelpath"`                                    // LayerLevelpath, Such As "1,2,3"
	Disabled        bool                  `json:"disabled" gorm:"-"`                                                                   // Disablestatus, According tostatusFieldDynamicCalculate
}

func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_details"
}
