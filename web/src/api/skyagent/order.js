import service from '@/utils/request'

// @Summary Get order list with pagination and filters
// @Router /admin-api/v1/orders [get]
export const getOrderList = (params) => {
  return service({
    url: '/admin-api/v1/orders',
    method: 'get',
    params
  })
}

// @Summary Get order detail by ID
// @Router /admin-api/v1/orders/:id [get]
export const getOrderDetail = (id) => {
  return service({
    url: `/admin-api/v1/orders/${id}`,
    method: 'get'
  })
}
