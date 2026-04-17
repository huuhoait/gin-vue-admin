package router

import (
	{{if .OnlyTemplate }} // {{end}}"{{.Module}}/middleware"
	"github.com/gin-gonic/gin"
)

var {{.StructName}} = new({{.Abbreviation}})

type {{.Abbreviation}} struct {}

// Init Initialize {{.Description}} router
func (r *{{.Abbreviation}}) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
{{- if not .OnlyTemplate }}
	{
	    group := private.Group("{{.Abbreviation}}").Use(middleware.OperationRecord())
		group.POST("create{{.StructName}}", api{{.StructName}}.Create{{.StructName}})   // Create {{.Description}}
		group.DELETE("delete{{.StructName}}", api{{.StructName}}.Delete{{.StructName}}) // Delete {{.Description}}
		group.DELETE("delete{{.StructName}}ByIds", api{{.StructName}}.Delete{{.StructName}}ByIds) // Batch delete {{.Description}}
		group.PUT("update{{.StructName}}", api{{.StructName}}.Update{{.StructName}})    // Update {{.Description}}
	}
	{
	    group := private.Group("{{.Abbreviation}}")
		group.GET("find{{.StructName}}", api{{.StructName}}.Find{{.StructName}})        // Get {{.Description}} by ID
		group.GET("get{{.StructName}}List", api{{.StructName}}.Get{{.StructName}}List)  // Get {{.Description}} list
	}
	{
	    group := public.Group("{{.Abbreviation}}")
    	{{- if .HasDataSource}}
	    group.GET("get{{.StructName}}DataSource", api{{.StructName}}.Get{{.StructName}}DataSource)  // Get {{.Description}} data source
	    {{- end}}
	    group.GET("get{{.StructName}}Public", api{{.StructName}}.Get{{.StructName}}Public)  // Public {{.Description}} endpoint
	}
{{- else}}
     // {
	 //   group := private.Group("{{.Abbreviation}}").Use(middleware.OperationRecord())
	 // }
	 // {
     //   group := private.Group("{{.Abbreviation}}")
     // }
    {
	    group := public.Group("{{.Abbreviation}}")
	    group.GET("get{{.StructName}}Public", api{{.StructName}}.Get{{.StructName}}Public)  // Public {{.Description}} endpoint
    }
{{- end}}
}
