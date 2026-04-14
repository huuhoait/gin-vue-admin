import service from '@/utils/request'

// @Tags Info
// @Summary Create announcement
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Info true "Create announcement"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Created"}"
// @Router /info/createInfo [post]
export const createInfo = (data) => {
  return service({
    url: '/info/createInfo',
    method: 'post',
    data
  })
}

// @Tags Info
// @Summary Delete announcement
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Info true "Delete announcement"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /info/deleteInfo [delete]
export const deleteInfo = (params) => {
  return service({
    url: '/info/deleteInfo',
    method: 'delete',
    params
  })
}

// @Tags Info
// @Summary Delete announcements in bulk
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Bulk delete announcements"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /info/deleteInfo [delete]
export const deleteInfoByIds = (params) => {
  return service({
    url: '/info/deleteInfoByIds',
    method: 'delete',
    params
  })
}

// @Tags Info
// @Summary Update announcement
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Info true "Update announcement"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /info/updateInfo [put]
export const updateInfo = (data) => {
  return service({
    url: '/info/updateInfo',
    method: 'put',
    data
  })
}

// @Tags Info
// @Summary Find announcement by id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Info true "Find announcement by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /info/findInfo [get]
export const findInfo = (params) => {
  return service({
    url: '/info/findInfo',
    method: 'get',
    params
  })
}

// @Tags Info
// @Summary Paginated announcement list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "Paginated announcement list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /info/getInfoList [get]
export const getInfoList = (params) => {
  return service({
    url: '/info/getInfoList',
    method: 'get',
    params
  })
}
// @Tags Info
// @Summary Get data source
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /info/findInfoDataSource [get]
export const getInfoDataSource = () => {
  return service({
    url: '/info/getInfoDataSource',
    method: 'get'
  })
}
