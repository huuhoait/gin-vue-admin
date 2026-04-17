package request

import (
	"encoding/json"
	"fmt"
	"github.com/huuhoait/gin-vue-admin/server/global"
	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/pkg/errors"
	"go/token"
	"strings"
)

type AutoCode struct {
	Package             string                 `json:"package"`
	PackageT            string                 `json:"-"`
	TableName           string                 `json:"tableName" example:"table name"`              // table name
	BusinessDB          string                 `json:"businessDB" example:"business database"`          // business database
	StructName          string                 `json:"structName" example:"Structname"`       // Structname
	PackageName         string                 `json:"packageName" example:"file name"`          // file name
	Description         string                 `json:"description" example:"StructChinese name"`    // StructChinese name
	Abbreviation        string                 `json:"abbreviation" example:"Structalias"`     // Structalias
	HumpPackageName     string                 `json:"humpPackageName" example:"gofile name"`    // gofile name
	GvaModel            bool                   `json:"gvaModel" example:"false"`            // whether to use gva default model
	AutoMigrate         bool                   `json:"autoMigrate" example:"false"`         // whether to auto-migrate schema
	AutoCreateResource  bool                   `json:"autoCreateResource" example:"false"`  // whether to auto-create resource id
	AutoCreateApiToSql  bool                   `json:"autoCreateApiToSql" example:"false"`  // whether to auto-create APIs
	AutoCreateMenuToSql bool                   `json:"autoCreateMenuToSql" example:"false"` // whether to auto-create menus
	AutoCreateBtnAuth   bool                   `json:"autoCreateBtnAuth" example:"false"`   // whether to auto-create button permissions
	OnlyTemplate        bool                   `json:"onlyTemplate" example:"false"`        // whether to generate template only
	IsTree              bool                   `json:"isTree" example:"false"`              // whether tree structure
	TreeJson            string                 `json:"treeJson" example:"JSON field for displayed tree"`       // JSON field for displayed tree
	IsAdd               bool                   `json:"isAdd" example:"false"`               // whether add operation
	Fields              []*AutoCodeField       `json:"fields"`
	GenerateWeb         bool                   `json:"generateWeb" example:"true"`    // whether generate web
	GenerateServer      bool                   `json:"generateServer" example:"true"` // whether generate server
	Module              string                 `json:"-"`
	DictTypes           []string               `json:"-"`
	PrimaryField        *AutoCodeField         `json:"primaryField"`
	DataSourceMap       map[string]*DataSource `json:"-"`
	HasPic              bool                   `json:"-"`
	HasFile             bool                   `json:"-"`
	HasTimer            bool                   `json:"-"`
	NeedSort            bool                   `json:"-"`
	NeedJSON            bool                   `json:"-"`
	HasRichText         bool                   `json:"-"`
	HasDataSource       bool                   `json:"-"`
	HasSearchTimer      bool                   `json:"-"`
	HasArray            bool                   `json:"-"`
	HasExcel            bool                   `json:"-"`
}

type DataSource struct {
	DBName       string `json:"dbName"`
	Table        string `json:"table"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Association  int    `json:"association"` // association: 1 one-to-one, 2 one-to-many
	HasDeletedAt bool   `json:"hasDeletedAt"`
}

func (r *AutoCode) Apis() []model.SysApi {
	return []model.SysApi{
		{
			Path:        "/" + r.Abbreviation + "/" + "create" + r.StructName,
			Description: "create" + r.Description,
			ApiGroup:    r.Description,
			Method:      "POST",
		},
		{
			Path:        "/" + r.Abbreviation + "/" + "delete" + r.StructName,
			Description: "delete" + r.Description,
			ApiGroup:    r.Description,
			Method:      "DELETE",
		},
		{
			Path:        "/" + r.Abbreviation + "/" + "delete" + r.StructName + "ByIds",
			Description: "batch delete" + r.Description,
			ApiGroup:    r.Description,
			Method:      "DELETE",
		},
		{
			Path:        "/" + r.Abbreviation + "/" + "update" + r.StructName,
			Description: "update" + r.Description,
			ApiGroup:    r.Description,
			Method:      "PUT",
		},
		{
			Path:        "/" + r.Abbreviation + "/" + "find" + r.StructName,
			Description: "get by ID" + r.Description,
			ApiGroup:    r.Description,
			Method:      "GET",
		},
		{
			Path:        "/" + r.Abbreviation + "/" + "get" + r.StructName + "List",
			Description: "get" + r.Description + "List",
			ApiGroup:    r.Description,
			Method:      "GET",
		},
	}
}

func (r *AutoCode) Menu(template string) model.SysBaseMenu {
	component := fmt.Sprintf("view/%s/%s/%s.vue", r.Package, r.PackageName, r.PackageName)
	if template != "package" {
		component = fmt.Sprintf("plugin/%s/view/%s.vue", r.Package, r.PackageName)
	}
	return model.SysBaseMenu{
		ParentId:  0,
		Path:      r.Abbreviation,
		Name:      r.Abbreviation,
		Component: component,
		Meta: model.Meta{
			Title: r.Description,
		},
	}
}

// Pretreatment preprocessing
// Author [SliverHorn](https://github.com/SliverHorn)
func (r *AutoCode) Pretreatment() error {
	r.Module = global.GVA_CONFIG.AutoCode.Module
	if token.IsKeyword(r.Abbreviation) {
		r.Abbreviation = r.Abbreviation + "_"
	} // Go keyword handling
	if strings.HasSuffix(r.HumpPackageName, "test") {
		r.HumpPackageName = r.HumpPackageName + "_"
	} // test
	length := len(r.Fields)
	dict := make(map[string]string, length)
	r.DataSourceMap = make(map[string]*DataSource, length)
	for i := 0; i < length; i++ {
		if r.Fields[i].Excel {
			r.HasExcel = true
		}
		if r.Fields[i].DictType != "" {
			dict[r.Fields[i].DictType] = ""
		}
		if r.Fields[i].Sort {
			r.NeedSort = true
		}
		switch r.Fields[i].FieldType {
		case "file":
			r.HasFile = true
			r.NeedJSON = true
		case "json":
			r.NeedJSON = true
		case "array":
			r.NeedJSON = true
			r.HasArray = true
		case "video":
			r.HasPic = true
		case "richtext":
			r.HasRichText = true
		case "picture":
			r.HasPic = true
		case "pictures":
			r.HasPic = true
			r.NeedJSON = true
		case "time.Time":
			r.HasTimer = true
			if r.Fields[i].FieldSearchType != "" && r.Fields[i].FieldSearchType != "BETWEEN" && r.Fields[i].FieldSearchType != "NOT BETWEEN" {
				r.HasSearchTimer = true
			}
		}
		if r.Fields[i].DataSource != nil {
			if r.Fields[i].DataSource.Table != "" && r.Fields[i].DataSource.Label != "" && r.Fields[i].DataSource.Value != "" {
				r.HasDataSource = true
				r.Fields[i].CheckDataSource = true
				r.DataSourceMap[r.Fields[i].FieldJson] = r.Fields[i].DataSource
			}
		}
		if !r.GvaModel && r.PrimaryField == nil && r.Fields[i].PrimaryKey {
			r.PrimaryField = r.Fields[i]
		} // custom primary key
	}
	{
		for key := range dict {
			r.DictTypes = append(r.DictTypes, key)
		}
	} // DictTypes => dictionary
	{
		if r.GvaModel {
			r.PrimaryField = &AutoCodeField{
				FieldName:    "ID",
				FieldType:    "uint",
				FieldDesc:    "ID",
				FieldJson:    "ID",
				DataTypeLong: "20",
				Comment:      "primary key ID",
				ColumnName:   "id",
			}
		}
	} // GvaModel
	{
		if r.IsAdd && r.PrimaryField == nil {
			r.PrimaryField = new(AutoCodeField)
		}
	} // ignore primary key in add-field mode
	if r.Package == "" {
		return errors.New("Package is empty!")
	} // extra check: Package not empty
	packages := []rune(r.Package)
	if len(packages) > 0 {
		if packages[0] >= 97 && packages[0] <= 122 {
			packages[0] = packages[0] - 32
		}
		r.PackageT = string(packages)
	} // PackageT: capitalized Package
	return nil
}

func (r *AutoCode) History() SysAutoHistoryCreate {
	bytes, _ := json.Marshal(r)
	return SysAutoHistoryCreate{
		Table:       r.TableName,
		Package:     r.Package,
		Request:     string(bytes),
		StructName:  r.StructName,
		BusinessDB:  r.BusinessDB,
		Description: r.Description,
	}
}

type AutoCodeField struct {
	FieldName       string `json:"fieldName"`       // FieldName
	FieldDesc       string `json:"fieldDesc"`       // InTextName
	FieldType       string `json:"fieldType"`       // FieldDatatype
	FieldJson       string `json:"fieldJson"`       // FieldJson
	DataTypeLong    string `json:"dataTypeLong"`    // DatabaseFieldLength
	Comment         string `json:"comment"`         // DatabaseFielddescription
	ColumnName      string `json:"columnName"`      // DatabaseField
	FieldSearchType string `json:"fieldSearchType"` // search conditions
	FieldSearchHide bool   `json:"fieldSearchHide"` // is hiddenQueryCondition
	DictType        string `json:"dictType"`        // Dictionary
	//Front           bool        `json:"front"`           // YesNoFrontendVisible
	Form            bool        `json:"form"`            // YesNoFrontendcreate/edit
	Table           bool        `json:"table"`           // YesNoFrontendTablegrid
	Desc            bool        `json:"desc"`            // YesNoFrontendDetails
	Excel           bool        `json:"excel"`           // YesNoimport/export
	Require         bool        `json:"require"`         // YesNoRequired
	DefaultValue    string      `json:"defaultValue"`    // YesNoRequired
	ErrorText       string      `json:"errorText"`       // validation error text
	Clearable       bool        `json:"clearable"`       // YesNoclearable
	Sort            bool        `json:"sort"`            // YesNoIncreasesort
	PrimaryKey      bool        `json:"primaryKey"`      // YesNoPrimary Key
	DataSource      *DataSource `json:"dataSource"`      // data source
	CheckDataSource bool        `json:"checkDataSource"` // YesNocheck data source
	FieldIndexType  string      `json:"fieldIndexType"`  // Indextype
}

type AutoFunc struct {
	Package         string `json:"package"`
	FuncName        string `json:"funcName"`        // methodname
	Router          string `json:"router"`          // route name
	FuncDesc        string `json:"funcDesc"`        // methodIntroduction
	BusinessDB      string `json:"businessDB"`      // business database
	StructName      string `json:"structName"`      // Structname
	PackageName     string `json:"packageName"`     // file name
	Description     string `json:"description"`     // StructChinese name
	Abbreviation    string `json:"abbreviation"`    // Structalias
	HumpPackageName string `json:"humpPackageName"` // gofile name
	Method          string `json:"method"`          // method
	IsPlugin        bool   `json:"isPlugin"`        // YesNoPlugin
	IsAuth          bool   `json:"isAuth"`          // YesNoAuthorization
	IsPreview       bool   `json:"isPreview"`       // YesNoPreview
	IsAi            bool   `json:"isAi"`            // YesNoAI
	ApiFunc         string `json:"apiFunc"`         // APImethod
	ServerFunc      string `json:"serverFunc"`      // Servicemethod
	JsFunc          string `json:"jsFunc"`          // JSmethod
}

type InitMenu struct {
	PlugName   string `json:"plugName"`
	ParentMenu string `json:"parentMenu"`
	Menus      []uint `json:"menus"`
}

type InitApi struct {
	PlugName string `json:"plugName"`
	APIs     []uint `json:"apis"`
}

type InitDictionary struct {
	PlugName     string `json:"plugName"`
	Dictionaries []uint `json:"dictionaries"`
}

type LLMAutoCode struct {
	Prompt string `json:"prompt" form:"prompt" gorm:"column:prompt;comment:PromptLanguage;type:text;"` //PromptLanguage
	Mode   string `json:"mode" form:"mode" gorm:"column:mode;comment:Mode;type:text;"`        //Mode
}
