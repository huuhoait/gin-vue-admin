package ast

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"os"
	"path/filepath"
	"testing"
)

func TestPluginInitializeGorm_Injection(t *testing.T) {
	type fields struct {
		Type        Type
		Path        string
		ImportPath  string
		StructName  string
		PackageName string
		IsNew       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test &model.User{} inject",
			fields: fields{
				Type:        TypePluginInitializeGorm,
				Path:        filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva", "initialize", "gorm.go"),
				ImportPath:  `"github.com/huuhoait/gin-vue-admin/server/plugin/gva/model"`,
				StructName:  "User",
				PackageName: "model",
				IsNew:       false,
			},
		},
		{
			name: "test new(model.ExaCustomer) inject",
			fields: fields{
				Type:        TypePluginInitializeGorm,
				Path:        filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva", "initialize", "gorm.go"),
				ImportPath:  `"github.com/huuhoait/gin-vue-admin/server/plugin/gva/model"`,
				StructName:  "User",
				PackageName: "model",
				IsNew:       true,
			},
		},
		{
			name: "test new(model.SysUsers) inject",
			fields: fields{
				Type:        TypePluginInitializeGorm,
				Path:        filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva", "initialize", "gorm.go"),
				ImportPath:  `"github.com/huuhoait/gin-vue-admin/server/plugin/gva/model"`,
				StructName:  "SysUser",
				PackageName: "model",
				IsNew:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(tt.fields.Path); err != nil {
				t.Skipf("skip: %s not present: %v", tt.fields.Path, err)
			}
			a := &PluginInitializeGorm{
				Type:        tt.fields.Type,
				Path:        tt.fields.Path,
				ImportPath:  tt.fields.ImportPath,
				StructName:  tt.fields.StructName,
				PackageName: tt.fields.PackageName,
				IsNew:       tt.fields.IsNew,
			}
			file, err := a.Parse(a.Path, nil)
			if err != nil {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if file != nil {
				a.Injection(file)
				err = a.Format(a.Path, nil, file)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Injection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginInitializeGorm_Rollback(t *testing.T) {
	type fields struct {
		Type        Type
		Path        string
		ImportPath  string
		StructName  string
		PackageName string
		IsNew       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test &model.User{} rollback",
			fields: fields{
				Type:        TypePluginInitializeGorm,
				Path:        filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva", "initialize", "gorm.go"),
				ImportPath:  `"github.com/huuhoait/gin-vue-admin/server/plugin/gva/model"`,
				StructName:  "User",
				PackageName: "model",
				IsNew:       false,
			},
		},
		{
			name: "test new(model.ExaCustomer) rollback",
			fields: fields{
				Type:        TypePluginInitializeGorm,
				Path:        filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva", "initialize", "gorm.go"),
				ImportPath:  `"github.com/huuhoait/gin-vue-admin/server/plugin/gva/model"`,
				StructName:  "User",
				PackageName: "model",
				IsNew:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(tt.fields.Path); err != nil {
				t.Skipf("skip: %s not present: %v", tt.fields.Path, err)
			}
			a := &PluginInitializeGorm{
				Type:        tt.fields.Type,
				Path:        tt.fields.Path,
				ImportPath:  tt.fields.ImportPath,
				StructName:  tt.fields.StructName,
				PackageName: tt.fields.PackageName,
				IsNew:       tt.fields.IsNew,
			}
			file, err := a.Parse(a.Path, nil)
			if err != nil {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if file != nil {
				a.Rollback(file)
				err = a.Format(a.Path, nil, file)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Rollback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
