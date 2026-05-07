import service from '@/utils/request'

// @Summary Get product list with pagination and filters
// @Router /admin-api/v1/products [get]
export const getProductList = (params) => {
  return service({
    url: '/admin-api/v1/products',
    method: 'get',
    params
  })
}

// @Summary Get supplier list
// @Router /admin-api/v1/suppliers [get]
export const getSupplierList = (params) => {
  return service({
    url: '/admin-api/v1/suppliers',
    method: 'get',
    params
  })
}
