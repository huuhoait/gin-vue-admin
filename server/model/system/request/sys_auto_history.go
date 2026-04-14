package request

import (
	common "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type SysAutoHistoryCreate struct {
	Table            string            // table name
	Package          string            // ModuleName/PluginName
	Request          string            // FrontendTransmitInofStructuretransformInformation
	StructName       string            // struct name
	BusinessDB       string            // business database
	Description      string            // StructChinese name
	Injections       map[string]string // injectpath
	Templates        map[string]string // template information
	ApiIDs           []uint            // apiTableRegistercontent
	MenuID           uint              // menu ID
	ExportTemplateID uint              // Export TemplateID
}

func (r *SysAutoHistoryCreate) Create() model.SysAutoCodeHistory {
	entity := model.SysAutoCodeHistory{
		Package:          r.Package,
		Request:          r.Request,
		Table:            r.Table,
		StructName:       r.StructName,
		Abbreviation:     r.StructName,
		BusinessDB:       r.BusinessDB,
		Description:      r.Description,
		Injections:       r.Injections,
		Templates:        r.Templates,
		ApiIDs:           r.ApiIDs,
		MenuID:           r.MenuID,
		ExportTemplateID: r.ExportTemplateID,
	}
	if entity.Table == "" {
		entity.Table = r.StructName
	}
	return entity
}

type SysAutoHistoryRollBack struct {
	common.GetById
	DeleteApi   bool `json:"deleteApi" form:"deleteApi"`     // YesNodeleteAPI
	DeleteMenu  bool `json:"deleteMenu" form:"deleteMenu"`   // YesNodelete menu
	DeleteTable bool `json:"deleteTable" form:"deleteTable"` // YesNodeleteTable
}

func (r *SysAutoHistoryRollBack) ApiIds(entity model.SysAutoCodeHistory) common.IdsReq {
	length := len(entity.ApiIDs)
	ids := make([]int, 0)
	for i := 0; i < length; i++ {
		ids = append(ids, int(entity.ApiIDs[i]))
	}
	return common.IdsReq{Ids: ids}
}
