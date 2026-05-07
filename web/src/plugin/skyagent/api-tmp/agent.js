import service from '@/utils/request'

// @Summary Get agent list with pagination and filters
// @Router /admin-api/v1/agents [get]
export const getAgentList = (params) => {
  return service({
    url: '/admin-api/v1/agents',
    method: 'get',
    params
  })
}

// @Summary Get agent detail by ID
// @Router /admin-api/v1/agents/:id [get]
export const getAgentDetail = (id) => {
  return service({
    url: `/admin-api/v1/agents/${id}`,
    method: 'get'
  })
}

// @Summary Get agent detail with plaintext PII (admin only)
// @Router /admin-api/v1/agents/:id/full [get]
// Upstream returns `Cache-Control: no-store` per contract §14.1.
export const getAgentAdminDetail = (id) => {
  return service({
    url: `/admin-api/v1/agents/${id}/full`,
    method: 'get'
  })
}

// @Summary Create a new agent
// @Router /admin-api/v1/agents [post]
export const createAgent = (data) => {
  return service({
    url: '/admin-api/v1/agents',
    method: 'post',
    data
  })
}

// @Summary Update agent info
// @Router /admin-api/v1/agents/:id [put]
export const updateAgent = (id, data) => {
  return service({
    url: `/admin-api/v1/agents/${id}`,
    method: 'put',
    data
  })
}

// @Summary Update agent status (approve/suspend/terminate)
// @Router /admin-api/v1/agents/:id/status [put]
export const updateAgentStatus = (id, data) => {
  return service({
    url: `/admin-api/v1/agents/${id}/status`,
    method: 'put',
    data
  })
}
