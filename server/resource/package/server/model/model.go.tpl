{{- if .IsAdd}}
// Add the following fields to the struct
{{- range .Fields}}
  {{ GenerateField . }}
{{- end }}

{{ else }}
// Auto-generated template {{.StructName}}
package {{.Package}}

{{- if not .OnlyTemplate}}
import (
	{{- if .GvaModel }}
	"{{.Module}}/global"
	{{- end }}
	{{- if or .HasTimer }}
	"time"
	{{- end }}
	{{- if .NeedJSON }}
	"gorm.io/datatypes"
	{{- end }}
)
{{- end }}

// {{.Description}} struct {{.StructName}}
type {{.StructName}} struct {
{{- if not .OnlyTemplate}}
{{- if .GvaModel }}
    global.GVA_MODEL
{{- end }}
{{- range .Fields}}
  {{ GenerateField . }}
{{- end }}
    {{- if .AutoCreateResource }}
    CreatedBy  uint   `gorm:"column:created_by;comment:Created by"`
    UpdatedBy  uint   `gorm:"column:updated_by;comment:Updated by"`
    DeletedBy  uint   `gorm:"column:deleted_by;comment:Deleted by"`
    {{- end }}
    {{- if .IsTree }}
    Children   []*{{.StructName}} `json:"children" gorm:"-"`     //Child nodes
    ParentID   int             `json:"parentID" gorm:"column:parent_id;comment:Parent node"`
    {{- end }}
{{- end }}
}

{{ if .TableName }}
// TableName {{.Description}} {{.StructName}} custom table name {{.TableName}}
func ({{.StructName}}) TableName() string {
    return "{{.TableName}}"
}
{{ end }}

{{if .IsTree }}
// GetChildren implements TreeNode
func (s *{{.StructName}}) GetChildren() []*{{.StructName}} {
    return s.Children
}

// SetChildren implements TreeNode
func (s *{{.StructName}}) SetChildren(children *{{.StructName}}) {
	s.Children = append(s.Children, children)
}

// GetID implements TreeNode
func (s *{{.StructName}}) GetID() int {
    return int({{if not .GvaModel}}*{{- end }}s.{{.PrimaryField.FieldName}})
}

// GetParentID implements TreeNode
func (s *{{.StructName}}) GetParentID() int {
    return s.ParentID
}
{{ end }}

{{ end }}
