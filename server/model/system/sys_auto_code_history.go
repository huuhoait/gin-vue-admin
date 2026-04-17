package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"gorm.io/gorm"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// SysAutoCodeHistory AutomaticmigrateMoveGenerationCodeRecord,Used forrollback,RepeatPlaceUse
type SysAutoCodeHistory struct {
	global.GVA_MODEL
	Table            string             `json:"tableName" gorm:"column:table_name;comment:table name"`
	Package          string             `json:"package" gorm:"column:package;comment:ModuleName/PluginName"`
	Request          string             `json:"request" gorm:"type:text;column:request;comment:FrontendTransmitInofStructuretransformInformation"`
	StructName       string             `json:"structName" gorm:"column:struct_name;comment:struct name"`
	Abbreviation     string             `json:"abbreviation" gorm:"column:abbreviation;comment:struct nameshrinkWrite"`
	BusinessDB       string             `json:"businessDb" gorm:"column:business_db;comment:business database"`
	Description      string             `json:"description" gorm:"column:description;comment:StructChinese name"`
	Templates        map[string]string  `json:"template" gorm:"serializer:json;type:text;column:templates;comment:template information"`
	Injections       map[string]string  `json:"injections" gorm:"serializer:json;type:text;column:Injections;comment:injectpath"`
	Flag             int                `json:"flag" gorm:"column:flag;comment:[0:Create,1:rollback]"`
	ApiIDs           []uint             `json:"apiIDs" gorm:"serializer:json;column:api_ids;comment:apiTableRegistercontent"`
	MenuID           uint               `json:"menuId" gorm:"column:menu_id;comment:menu ID"`
	ExportTemplateID uint               `json:"exportTemplateID" gorm:"column:export_template_id;comment:Export TemplateID"`
	AutoCodePackage  SysAutoCodePackage `json:"autoCodePackage" gorm:"foreignKey:ID;references:PackageID"`
	PackageID        uint               `json:"packageID" gorm:"column:package_id;comment:PackageID"`
}

func (s *SysAutoCodeHistory) BeforeCreate(db *gorm.DB) error {
	templates := make(map[string]string, len(s.Templates))
	for key, value := range s.Templates {
		server := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server)
		{
			hasServer := strings.Index(key, server)
			if hasServer != -1 {
				key = strings.TrimPrefix(key, server)
				keys := strings.Split(key, string(os.PathSeparator))
				key = path.Join(keys...)
			}
		} // key
		web := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.WebRoot())
		hasWeb := strings.Index(value, web)
		if hasWeb != -1 {
			value = strings.TrimPrefix(value, web)
			values := strings.Split(value, string(os.PathSeparator))
			value = path.Join(values...)
			templates[key] = value
			continue
		}
		hasServer := strings.Index(value, server)
		if hasServer != -1 {
			value = strings.TrimPrefix(value, server)
			values := strings.Split(value, string(os.PathSeparator))
			value = path.Join(values...)
			templates[key] = value
			continue
		}
	}
	s.Templates = templates
	return nil
}

func (s *SysAutoCodeHistory) TableName() string {
	return "sys_auto_code_histories"
}
