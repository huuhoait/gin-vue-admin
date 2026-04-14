import service from '@/utils/request'
// @Tags SysDictionaryDetail
// @Summary Create SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionaryDetail true "Create SysDictionaryDetail"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/createSysDictionaryDetail [post]
export const createSysDictionaryDetail = (data) => {
  return service({
    url: '/sysDictionaryDetail/createSysDictionaryDetail',
    method: 'post',
    data
  })
}

// @Tags SysDictionaryDetail
// @Summary Delete SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionaryDetail true "Delete SysDictionaryDetail"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysDictionaryDetail/deleteSysDictionaryDetail [delete]
export const deleteSysDictionaryDetail = (data) => {
  return service({
    url: '/sysDictionaryDetail/deleteSysDictionaryDetail',
    method: 'delete',
    data
  })
}

// @Tags SysDictionaryDetail
// @Summary Update SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionaryDetail true "Update SysDictionaryDetail"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /sysDictionaryDetail/updateSysDictionaryDetail [put]
export const updateSysDictionaryDetail = (data) => {
  return service({
    url: '/sysDictionaryDetail/updateSysDictionaryDetail',
    method: 'put',
    data
  })
}

// @Tags SysDictionaryDetail
// @Summary Find SysDictionaryDetail by id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysDictionaryDetail true "Find SysDictionaryDetail by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/findSysDictionaryDetail [get]
export const findSysDictionaryDetail = (params) => {
  return service({
    url: '/sysDictionaryDetail/findSysDictionaryDetail',
    method: 'get',
    params
  })
}

// @Tags SysDictionaryDetail
// @Summary Paginated SysDictionaryDetail list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "Paginated SysDictionaryDetail list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/getSysDictionaryDetailList [get]
export const getSysDictionaryDetailList = (params) => {
  return service({
    url: '/sysDictionaryDetail/getSysDictionaryDetailList',
    method: 'get',
    params
  })
}

// @Tags SysDictionaryDetail
// @Summary Get dictionary detail tree (by dictionary ID)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param sysDictionaryID query string true "Dictionary ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/getDictionaryTreeList [get]
export const getDictionaryTreeList = (params) => {
  return service({
    url: '/sysDictionaryDetail/getDictionaryTreeList',
    method: 'get',
    params
  })
}

// @Tags SysDictionaryDetail
// @Summary Get dictionary detail tree (by dictionary type)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param dictType query string true "Dictionary type"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/getDictionaryTreeListByType [get]
export const getDictionaryTreeListByType = (params) => {
  return service({
    url: '/sysDictionaryDetail/getDictionaryTreeListByType',
    method: 'get',
    params
  })
}

// @Tags SysDictionaryDetail
// @Summary Get dictionary details by parent ID
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param parentID query string true "Parent ID"
// @Param includeChildren query boolean false "Include children"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/getDictionaryDetailsByParent [get]
export const getDictionaryDetailsByParent = (params) => {
  return service({
    url: '/sysDictionaryDetail/getDictionaryDetailsByParent',
    method: 'get',
    params
  })
}

// @Tags SysDictionaryDetail
// @Summary Get full path of a dictionary detail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query string true "Dictionary detail ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysDictionaryDetail/getDictionaryPath [get]
export const getDictionaryPath = (params) => {
  return service({
    url: '/sysDictionaryDetail/getDictionaryPath',
    method: 'get',
    params
  })
}
