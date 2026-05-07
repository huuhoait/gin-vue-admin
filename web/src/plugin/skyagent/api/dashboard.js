import service from '@/utils/request'

// @Summary Get dashboard overview metrics
// @Router /admin-api/v1/dashboard/overview [get]
export const getDashboardOverview = () => {
  return service({
    url: '/admin-api/v1/dashboard/overview',
    method: 'get'
  })
}
