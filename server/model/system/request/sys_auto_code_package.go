package request

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	model "github.com/huuhoait/gin-vue-admin/server/model/system"
)

type SysAutoCodePackageCreate struct {
	Desc        string `json:"desc" example:"Description"`
	Label       string `json:"label" example:"Display name"`
	Template    string `json:"template"  example:"Template"`
	PackageName string `json:"packageName" example:"Package name"`
	Module      string `json:"-" example:"Module"`
}

func (r *SysAutoCodePackageCreate) AutoCode() AutoCode {
	return AutoCode{
		Package: r.PackageName,
		Module:  global.GVA_CONFIG.AutoCode.Module,
	}
}

func (r *SysAutoCodePackageCreate) Create() model.SysAutoCodePackage {
	return model.SysAutoCodePackage{
		Desc:        r.Desc,
		Label:       r.Label,
		Template:    r.Template,
		PackageName: r.PackageName,
		Module:      global.GVA_CONFIG.AutoCode.Module,
	}
}
