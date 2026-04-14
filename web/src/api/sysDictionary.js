import service from '@/utils/request'
// @Tags SysDictionary
// @Summary Create SysDictionary
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionary true "Create SysDictionary"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionary/createSysDictionary [post]
export const createSysDictionary = (data) => {
  return service({
    url: '/sysDictionary/createSysDictionary',
    method: 'post',
    data
  })
}

// @Tags SysDictionary
// @Summary Delete SysDictionary
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionary true "Delete SysDictionary"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysDictionary/deleteSysDictionary [delete]
export const deleteSysDictionary = (data) => {
  return service({
    url: '/sysDictionary/deleteSysDictionary',
    method: 'delete',
    data
  })
}

// @Tags SysDictionary
// @Summary Update SysDictionary
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionary true "Update SysDictionary"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /sysDictionary/updateSysDictionary [put]
export const updateSysDictionary = (data) => {
  return service({
    url: '/sysDictionary/updateSysDictionary',
    method: 'put',
    data
  })
}

// @Tags SysDictionary
// @Summary Find SysDictionary by id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionary true "Find SysDictionary by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionary/findSysDictionary [get]
export const findSysDictionary = (params) => {
  return service({
    url: '/sysDictionary/findSysDictionary',
    method: 'get',
    params
  })
}

// @Tags SysDictionary
// @Summary Paginated SysDictionary list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "Paginated SysDictionary list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionary/getSysDictionaryList [get]
export const getSysDictionaryList = (params) => {
  return service({
    url: '/sysDictionary/getSysDictionaryList',
    method: 'get',
    params
  })
}

// @Tags SysDictionary
// @Summary Export dictionary JSON (including details)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SysDictionary true "Dictionary ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Exported"}"
// @Router /sysDictionary/exportSysDictionary [get]
export const exportSysDictionary = (params) => {
  return service({
    url: '/sysDictionary/exportSysDictionary',
    method: 'get',
    params
  })
}

// @Tags SysDictionary
// @Summary Import dictionary JSON (including details)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "Dictionary JSON"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Imported"}"
// @Router /sysDictionary/importSysDictionary [post]
export const importSysDictionary = (data) => {
  return service({
    url: '/sysDictionary/importSysDictionary',
    method: 'post',
    data
  })
}
