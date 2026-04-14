import service from '@/utils/request'
// @Tags SysApi
// @Summary Create customer
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "Create customer"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /customer/customer [post]
export const createExaCustomer = (data) => {
  return service({
    url: '/customer/customer',
    method: 'post',
    data
  })
}

// @Tags SysApi
// @Summary Update customer
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "Update customer"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /customer/customer [put]
export const updateExaCustomer = (data) => {
  return service({
    url: '/customer/customer',
    method: 'put',
    data
  })
}

// @Tags SysApi
// @Summary Delete customer
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "Delete customer"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /customer/customer [delete]
export const deleteExaCustomer = (data) => {
  return service({
    url: '/customer/customer',
    method: 'delete',
    data
  })
}

// @Tags SysApi
// @Summary Get customer
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "Get customer"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /customer/customer [get]
export const getExaCustomer = (params) => {
  return service({
    url: '/customer/customer',
    method: 'get',
    params
  })
}

// @Tags SysApi
// @Summary Paginated customer list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "Paginated customer list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /customer/customerList [get]
export const getExaCustomerList = (params) => {
  return service({
    url: '/customer/customerList',
    method: 'get',
    params
  })
}
