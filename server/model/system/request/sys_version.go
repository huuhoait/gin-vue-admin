package request

import (
	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"time"
)

type SysVersionSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	VersionName    *string     `json:"versionName" form:"versionName"`
	VersionCode    *string     `json:"versionCode" form:"versionCode"`
	request.PageInfo
}

// ExportVersionRequest export versionRequestStructureBody
type ExportVersionRequest struct {
	VersionName string `json:"versionName" binding:"required"` // version name
	VersionCode string `json:"versionCode" binding:"required"` // version number
	Description string `json:"description"`                    // version description
	MenuIds     []uint `json:"menuIds"`                        // ChooseInmenu IDList
	ApiIds      []uint `json:"apiIds"`                         // ChooseInAPI IDList
	DictIds     []uint `json:"dictIds"`                        // ChooseInDictionaryIDList
}

// ImportVersionRequest import versionRequestStructureBody
type ImportVersionRequest struct {
	VersionInfo      VersionInfo            `json:"version" binding:"required"` // version information
	ExportMenu       []system.SysBaseMenu   `json:"menus"`                      // Menudata; reuse directlySysBaseMenu
	ExportApi        []system.SysApi        `json:"apis"`                       // APIdata; reuse directlySysApi
	ExportDictionary []system.SysDictionary `json:"dictionaries"`               // Dictionarydata; reuse directlySysDictionary
}

// VersionInfo version informationStructureBody
type VersionInfo struct {
	Name        string `json:"name" binding:"required"`        // version name
	Code        string `json:"code" binding:"required"`        // version number
	Description string `json:"description"`                    // version description
	ExportTime  string `json:"exportTime"`                     // exportTime
}
