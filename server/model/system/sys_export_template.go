// auto-generate templateSysExportTemplate
package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
)

// Export Template StructureBody  SysExportTemplate
type SysExportTemplate struct {
	global.GVA_MODEL
	DBName       string         `json:"dbName" form:"dbName" gorm:"column:db_name;comment:Databasename;"`                       //Databasename
	Name         string         `json:"name" form:"name" gorm:"column:name;comment:Templatename;"`                               //Templatename
	TableName    string         `json:"tableName" form:"tableName" gorm:"column:table_name;comment:table nameName;"`                //table nameName
	TemplateID   string         `json:"templateID" form:"templateID" gorm:"column:template_id;comment:template identifier;"`            //template identifier
	TemplateInfo string         `json:"templateInfo" form:"templateInfo" gorm:"column:template_info;type:text;"`         //template information
	SQL          string         `json:"sql" form:"sql" gorm:"column:sql;type:text;comment:CustomexportSQL;"`                    //CustomexportSQL
	ImportSQL    string         `json:"importSql" form:"importSql" gorm:"column:import_sql;type:text;comment:CustomimportSQL;"` //CustomimportSQL
	Limit        *int           `json:"limit" form:"limit" gorm:"column:limit;comment:exportLimit"`
	Order        string         `json:"order" form:"order" gorm:"column:order;comment:sort"`
	Conditions   []Condition    `json:"conditions" form:"conditions" gorm:"foreignKey:TemplateID;references:TemplateID;comment:Condition"`
	JoinTemplate []JoinTemplate `json:"joinTemplate" form:"joinTemplate" gorm:"foreignKey:TemplateID;references:TemplateID;comment:Association"`
}

type JoinTemplate struct {
	global.GVA_MODEL
	TemplateID string `json:"templateID" form:"templateID" gorm:"column:template_id;comment:template identifier"`
	JOINS      string `json:"joins" form:"joins" gorm:"column:joins;comment:Association"`
	Table      string `json:"table" form:"table" gorm:"column:table;comment:related table"`
	ON         string `json:"on" form:"on" gorm:"column:on;comment:AssociationCondition"`
}

func (JoinTemplate) TableName() string {
	return "sys_export_template_join"
}

type Condition struct {
	global.GVA_MODEL
	TemplateID string `json:"templateID" form:"templateID" gorm:"column:template_id;comment:template identifier"`
	From       string `json:"from" form:"from" gorm:"column:from;comment:ConditionFetchofkey"`
	Column     string `json:"column" form:"column" gorm:"column:column;comment:AsQueryConditionofField"`
	Operator   string `json:"operator" form:"operator" gorm:"column:operator;comment:OperationSymbol"`
}

func (Condition) TableName() string {
	return "sys_export_template_condition"
}
