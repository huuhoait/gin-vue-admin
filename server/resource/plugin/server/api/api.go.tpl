package api

import (
{{if not .OnlyTemplate}}
	"{{.Module}}/global"
    "{{.Module}}/model/common/response"
    "{{.Module}}/plugin/{{.Package}}/model"
    {{- if not .IsTree}}
    "{{.Module}}/plugin/{{.Package}}/model/request"
    {{- end }}
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    {{- if .AutoCreateResource}}
    "{{.Module}}/utils"
    {{- end }}
{{- else }}
    "{{.Module}}/model/common/response"
    "github.com/gin-gonic/gin"
{{- end }}
)

var {{.StructName}} = new({{.Abbreviation}})

type {{.Abbreviation}} struct {}
{{if not .OnlyTemplate}}
// Create{{.StructName}} Create {{.Description}}
// @Tags {{.StructName}}
// @Summary Create {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "Create {{.Description}}"
// @Success 200 {object} response.Response{msg=string} "Created successfully"
// @Router /{{.Abbreviation}}/create{{.StructName}} [post]
func (a *{{.Abbreviation}}) Create{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	var info model.{{.StructName}}
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	{{- if .AutoCreateResource }}
    info.CreatedBy = utils.GetUserID(c)
	{{- end }}
	err = service{{ .StructName }}.Create{{.StructName}}(ctx,&info)
	if err != nil {
        global.GVA_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed:" + err.Error(), c)
		return
	}
    response.OkWithMessage("Created successfully", c)
}

// Delete{{.StructName}} Delete {{.Description}}
// @Tags {{.StructName}}
// @Summary Delete {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "Delete {{.Description}}"
// @Success 200 {object} response.Response{msg=string} "Deleted successfully"
// @Router /{{.Abbreviation}}/delete{{.StructName}} [delete]
func (a *{{.Abbreviation}}) Delete{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	{{.PrimaryField.FieldJson}} := c.Query("{{.PrimaryField.FieldJson}}")
{{- if .AutoCreateResource }}
    userID := utils.GetUserID(c)
{{- end }}
	err := service{{ .StructName }}.Delete{{.StructName}}(ctx,{{.PrimaryField.FieldJson}} {{- if .AutoCreateResource -}},userID{{- end -}})
	if err != nil {
        global.GVA_LOG.Error("Deletion failed!", zap.Error(err))
		response.FailWithMessage("Deletion failed:" + err.Error(), c)
		return
	}
    response.OkWithMessage("Deleted successfully", c)
}

// Delete{{.StructName}}ByIds Batch delete {{.Description}}
// @Tags {{.StructName}}
// @Summary Batch delete {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "Batch deletion succeeded"
// @Router /{{.Abbreviation}}/delete{{.StructName}}ByIds [delete]
func (a *{{.Abbreviation}}) Delete{{.StructName}}ByIds(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	{{.PrimaryField.FieldJson}}s := c.QueryArray("{{.PrimaryField.FieldJson}}s[]")
{{- if .AutoCreateResource }}
    userID := utils.GetUserID(c)
{{- end }}
	err := service{{ .StructName }}.Delete{{.StructName}}ByIds(ctx,{{.PrimaryField.FieldJson}}s{{- if .AutoCreateResource }},userID{{- end }})
	if err != nil {
        global.GVA_LOG.Error("Batch deletion failed!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed:" + err.Error(), c)
		return
	}
    response.OkWithMessage("Batch deletion succeeded", c)
}

// Update{{.StructName}} Update {{.Description}}
// @Tags {{.StructName}}
// @Summary Update {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "Update {{.Description}}"
// @Success 200 {object} response.Response{msg=string} "Updated successfully"
// @Router /{{.Abbreviation}}/update{{.StructName}} [put]
func (a *{{.Abbreviation}}) Update{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	var info model.{{.StructName}}
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
{{- if .AutoCreateResource }}
    info.UpdatedBy = utils.GetUserID(c)
{{- end }}
	err = service{{ .StructName }}.Update{{.StructName}}(ctx,info)
    if err != nil {
        global.GVA_LOG.Error("Update failed!", zap.Error(err))
		response.FailWithMessage("Update failed:" + err.Error(), c)
		return
	}
    response.OkWithMessage("Updated successfully", c)
}

// Find{{.StructName}} Query {{.Description}} by id
// @Tags {{.StructName}}
// @Summary Query {{.Description}} by id
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param {{.PrimaryField.FieldJson}} query {{.PrimaryField.FieldType}} true "Query {{.Description}} by id"
// @Success 200 {object} response.Response{data=model.{{.StructName}},msg=string} "Query succeeded"
// @Router /{{.Abbreviation}}/find{{.StructName}} [get]
func (a *{{.Abbreviation}}) Find{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	{{.PrimaryField.FieldJson}} := c.Query("{{.PrimaryField.FieldJson}}")
	re{{.Abbreviation}}, err := service{{ .StructName }}.Get{{.StructName}}(ctx,{{.PrimaryField.FieldJson}})
	if err != nil {
        global.GVA_LOG.Error("Query failed!", zap.Error(err))
		response.FailWithMessage("Query failed:" + err.Error(), c)
		return
	}
    response.OkWithData(re{{.Abbreviation}}, c)
}

{{- if .IsTree }}
// Get{{.StructName}}List Paginated {{.Description}} list
// @Tags {{.StructName}}
// @Summary Paginated {{.Description}} list
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
func (a *{{.Abbreviation}}) Get{{.StructName}}List(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	list, err := service{{ .StructName }}.Get{{.StructName}}InfoList(ctx)
	if err != nil {
	    global.GVA_LOG.Error("Retrieval failed!", zap.Error(err))
        response.FailWithMessage("Retrieval failed:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(list, "Retrieved successfully", c)
}
{{- else }}
// Get{{.StructName}}List Paginated {{.Description}} list
// @Tags {{.StructName}}
// @Summary Paginated {{.Description}} list
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.{{.StructName}}Search true "Paginated {{.Description}} list"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
func (a *{{.Abbreviation}}) Get{{.StructName}}List(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	var pageInfo request.{{.StructName}}Search
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service{{ .StructName }}.Get{{.StructName}}InfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("Retrieval failed!", zap.Error(err))
        response.FailWithMessage("Retrieval failed:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "Retrieved successfully", c)
}
{{- end }}

{{- if .HasDataSource }}
// Get{{.StructName}}DataSource Get the data source for {{.StructName}}
// @Tags {{.StructName}}
// @Summary Get the data source for {{.StructName}}
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Query succeeded"
// @Router /{{.Abbreviation}}/get{{.StructName}}DataSource [get]
func (a *{{.Abbreviation}}) Get{{.StructName}}DataSource(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

    // This endpoint returns the data defined for the data source
   dataSource, err := service{{ .StructName }}.Get{{.StructName}}DataSource(ctx)
   if err != nil {
		global.GVA_LOG.Error("Query failed!", zap.Error(err))
        response.FailWithMessage("Query failed:" + err.Error(), c)
		return
   }
    response.OkWithData(dataSource, c)
}
{{- end }}
{{- end }}
// Get{{.StructName}}Public Public {{.Description}} endpoint (no auth required)
// @Tags {{.StructName}}
// @Summary Public {{.Description}} endpoint (no auth required)
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/get{{.StructName}}Public [get]
func (a *{{.Abbreviation}}) Get{{.StructName}}Public(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

    // This endpoint does not require authentication. The example returns a fixed message; typically this is used for public (C-side) services, and you should implement your own business logic
    service{{ .StructName }}.Get{{.StructName}}Public(ctx)
    response.OkWithDetailed(gin.H{"info": "Public {{.Description}} endpoint info (no auth required)"}, "Retrieved successfully", c)
}
