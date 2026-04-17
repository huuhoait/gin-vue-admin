package {{.Package}}

import (
	{{if not .OnlyTemplate}}
	"{{.Module}}/global"
    "{{.Module}}/model/common/response"
    "{{.Module}}/model/{{.Package}}"
    {{- if not .IsTree}}
    {{.Package}}Req "{{.Module}}/model/{{.Package}}/request"
    {{- end }}
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    {{- if .AutoCreateResource}}
    "{{.Module}}/utils"
    {{- end }}
    {{- else}}
    "{{.Module}}/model/common/response"
    "github.com/gin-gonic/gin"
    {{- end}}
)

type {{.StructName}}Api struct {}

{{if not .OnlyTemplate}}

// Create{{.StructName}} Create {{.Description}}
// @Tags {{.StructName}}
// @Summary Create {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body {{.Package}}.{{.StructName}} true "Create {{.Description}}"
// @Success 200 {object} response.Response{msg=string} "Created successfully"
// @Router /{{.Abbreviation}}/create{{.StructName}} [post]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Create{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	var {{.Abbreviation}} {{.Package}}.{{.StructName}}
	err := c.ShouldBindJSON(&{{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	{{- if .AutoCreateResource }}
    {{.Abbreviation}}.CreatedBy = utils.GetUserID(c)
	{{- end }}
	err = {{.Abbreviation}}Service.Create{{.StructName}}(ctx,&{{.Abbreviation}})
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
// @Param data body {{.Package}}.{{.StructName}} true "Delete {{.Description}}"
// @Success 200 {object} response.Response{msg=string} "Deleted successfully"
// @Router /{{.Abbreviation}}/delete{{.StructName}} [delete]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Delete{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	{{.PrimaryField.FieldJson}} := c.Query("{{.PrimaryField.FieldJson}}")
		{{- if .AutoCreateResource }}
    userID := utils.GetUserID(c)
        {{- end }}
	err := {{.Abbreviation}}Service.Delete{{.StructName}}(ctx,{{.PrimaryField.FieldJson}} {{- if .AutoCreateResource -}},userID{{- end -}})
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
func ({{.Abbreviation}}Api *{{.StructName}}Api) Delete{{.StructName}}ByIds(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	{{.PrimaryField.FieldJson}}s := c.QueryArray("{{.PrimaryField.FieldJson}}s[]")
    	{{- if .AutoCreateResource }}
    userID := utils.GetUserID(c)
        {{- end }}
	err := {{.Abbreviation}}Service.Delete{{.StructName}}ByIds(ctx,{{.PrimaryField.FieldJson}}s{{- if .AutoCreateResource }},userID{{- end }})
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
// @Param data body {{.Package}}.{{.StructName}} true "Update {{.Description}}"
// @Success 200 {object} response.Response{msg=string} "Updated successfully"
// @Router /{{.Abbreviation}}/update{{.StructName}} [put]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Update{{.StructName}}(c *gin.Context) {
    // Retrieve the standard context from ctx for the business operation
    ctx := c.Request.Context()

	var {{.Abbreviation}} {{.Package}}.{{.StructName}}
	err := c.ShouldBindJSON(&{{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	    {{- if .AutoCreateResource }}
    {{.Abbreviation}}.UpdatedBy = utils.GetUserID(c)
        {{- end }}
	err = {{.Abbreviation}}Service.Update{{.StructName}}(ctx,{{.Abbreviation}})
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
// @Success 200 {object} response.Response{data={{.Package}}.{{.StructName}},msg=string} "Query succeeded"
// @Router /{{.Abbreviation}}/find{{.StructName}} [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Find{{.StructName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	{{.PrimaryField.FieldJson}} := c.Query("{{.PrimaryField.FieldJson}}")
	re{{.Abbreviation}}, err := {{.Abbreviation}}Service.Get{{.StructName}}(ctx,{{.PrimaryField.FieldJson}})
	if err != nil {
        global.GVA_LOG.Error("Query failed!", zap.Error(err))
		response.FailWithMessage("Query failed:" + err.Error(), c)
		return
	}
	response.OkWithData(re{{.Abbreviation}}, c)
}

{{- if .IsTree }}
// Get{{.StructName}}List Paginated {{.Description}} list. Tree mode does not accept parameters
// @Tags {{.StructName}}
// @Summary Paginated {{.Description}} list
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}List(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	list, err := {{.Abbreviation}}Service.Get{{.StructName}}InfoList(ctx)
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
// @Param data query {{.Package}}Req.{{.StructName}}Search true "Paginated {{.Description}} list"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}List(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

	var pageInfo {{.Package}}Req.{{.StructName}}Search
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := {{.Abbreviation}}Service.Get{{.StructName}}InfoList(ctx,pageInfo)
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
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}DataSource(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

    // This endpoint returns the data defined for the data source
    dataSource, err := {{.Abbreviation}}Service.Get{{.StructName}}DataSource(ctx)
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
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}Public(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()

    // This endpoint does not require authentication
    // The example returns a fixed message. Typically this is used for public (C-side) services, and you should implement your own business logic
    {{.Abbreviation}}Service.Get{{.StructName}}Public(ctx)
    response.OkWithDetailed(gin.H{
       "info": "Public {{.Description}} endpoint info (no auth required)",
    }, "Retrieved successfully", c)
}
