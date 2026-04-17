{{if .IsPlugin}}
// {{.FuncName}} {{.FuncDesc}}
// @Tags {{.StructName}}
// @Summary {{.FuncDesc}}
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/{{.Router}} [{{.Method}}]
export const {{.Router}} = () => {
  return service({
    url: '/{{.Abbreviation}}/{{.Router}}',
    method: '{{.Method}}'
  })
}

{{- else -}}

// {{.FuncName}} {{.FuncDesc}}
// @Tags {{.StructName}}
// @Summary {{.FuncDesc}}
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Success"
// @Router /{{.Abbreviation}}/{{.Router}} [{{.Method}}]
export const {{.Router}} = () => {
  return service({
    url: '/{{.Abbreviation}}/{{.Router}}',
    method: '{{.Method}}'
  })
}

{{- end -}}
