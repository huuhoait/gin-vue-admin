import service from '@/utils/request'
{{- if not .OnlyTemplate}}
// @Tags {{.StructName}}
// @Summary Create {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "Create {{.Description}}"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Created successfully"}"
// @Router /{{.Abbreviation}}/create{{.StructName}} [post]
export const create{{.StructName}} = (data) => {
  return service({
    url: '/{{.Abbreviation}}/create{{.StructName}}',
    method: 'post',
    data
  })
}

// @Tags {{.StructName}}
// @Summary Delete {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "Delete {{.Description}}"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted successfully"}"
// @Router /{{.Abbreviation}}/delete{{.StructName}} [delete]
export const delete{{.StructName}} = (params) => {
  return service({
    url: '/{{.Abbreviation}}/delete{{.StructName}}',
    method: 'delete',
    params
  })
}

// @Tags {{.StructName}}
// @Summary Batch delete {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Batch delete {{.Description}}"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted successfully"}"
// @Router /{{.Abbreviation}}/delete{{.StructName}} [delete]
export const delete{{.StructName}}ByIds = (params) => {
  return service({
    url: '/{{.Abbreviation}}/delete{{.StructName}}ByIds',
    method: 'delete',
    params
  })
}

// @Tags {{.StructName}}
// @Summary Update {{.Description}}
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "Update {{.Description}}"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated successfully"}"
// @Router /{{.Abbreviation}}/update{{.StructName}} [put]
export const update{{.StructName}} = (data) => {
  return service({
    url: '/{{.Abbreviation}}/update{{.StructName}}',
    method: 'put',
    data
  })
}

// @Tags {{.StructName}}
// @Summary Query {{.Description}} by id
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.{{.StructName}} true "Query {{.Description}} by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Query succeeded"}"
// @Router /{{.Abbreviation}}/find{{.StructName}} [get]
export const find{{.StructName}} = (params) => {
  return service({
    url: '/{{.Abbreviation}}/find{{.StructName}}',
    method: 'get',
    params
  })
}

// @Tags {{.StructName}}
// @Summary Paginated {{.Description}} list
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "Paginated {{.Description}} list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Retrieved successfully"}"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
export const get{{.StructName}}List = (params) => {
  return service({
    url: '/{{.Abbreviation}}/get{{.StructName}}List',
    method: 'get',
    params
  })
}

{{- if .HasDataSource}}
// @Tags {{.StructName}}
// @Summary Get data source
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Query succeeded"}"
// @Router /{{.Abbreviation}}/find{{.StructName}}DataSource [get]
export const get{{.StructName}}DataSource = () => {
  return service({
    url: '/{{.Abbreviation}}/get{{.StructName}}DataSource',
    method: 'get',
  })
}
{{- end}}
{{- end}}
// @Tags {{.StructName}}
// @Summary Public {{.Description}} endpoint (no auth required)
// @Accept application/json
// @Produce application/json
// @Param data query request.{{.StructName}}Search true "Paginated {{.Description}} list"
// @Success 200 {object} response.Response{data=object,msg=string} "Retrieved successfully"
// @Router /{{.Abbreviation}}/get{{.StructName}}Public [get]
export const get{{.StructName}}Public = () => {
  return service({
    url: '/{{.Abbreviation}}/get{{.StructName}}Public',
    method: 'get',
  })
}
