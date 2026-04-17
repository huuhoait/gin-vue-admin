package response

import (
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
)

// ExportVersionResponse export versionResponseStructureBody
type ExportVersionResponse struct {
	Version      request.VersionInfo    `json:"version"`      // version information
	Menus        []system.SysBaseMenu   `json:"menus"`        // Menudata; reuse directlySysBaseMenu
	Apis         []system.SysApi        `json:"apis"`         // APIdata; reuse directlySysApi
	Dictionaries []system.SysDictionary `json:"dictionaries"` // Dictionarydata; reuse directlySysDictionary
}
