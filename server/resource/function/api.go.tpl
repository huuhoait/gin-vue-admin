{{if .IsPlugin}}
// {{.FuncName}} {{.FuncDesc}}
// @Tags {{.StructName}}
// @Summary {{.FuncDesc}}
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/{{.Router}} [{{.Method}}]
func (a *{{.Abbreviation}}) {{.FuncName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()
    // Add your business logic here
    err := service{{ .StructName }}.{{.FuncName}}(ctx)
       if err != nil {
    		global.GVA_LOG.Error("Failed!", zap.Error(err))
            response.FailWithMessage("Failed", c)
    		return
       }
    response.OkWithData("Returned data",c)
}

{{- else -}}

// {{.FuncName}} {{.FuncDesc}}
// @Tags {{.StructName}}
// @Summary {{.FuncDesc}}
// @Accept application/json
// @Produce application/json
// @Param data query {{.Package}}Req.{{.StructName}}Search true "Success"
// @Success 200 {object} response.Response{data=object,msg=string} "Success"
// @Router /{{.Abbreviation}}/{{.Router}} [{{.Method}}]
func ({{.Abbreviation}}Api *{{.StructName}}Api){{.FuncName}}(c *gin.Context) {
    // Create a business context
    ctx := c.Request.Context()
    // Add your business logic here
    err := {{.Abbreviation}}Service.{{.FuncName}}(ctx)
    if err != nil {
        global.GVA_LOG.Error("Failed!", zap.Error(err))
   		response.FailWithMessage("Failed", c)
   		return
   	}
   	response.OkWithData("Returned data",c)
}
{{end}}
