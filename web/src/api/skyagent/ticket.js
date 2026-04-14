import service from '@/utils/request'

export const createTicket = (data) => {
  return service({ url: '/admin-api/v1/onboarding/tickets', method: 'post', data })
}

export const getTicketList = (params) => {
  return service({ url: '/admin-api/v1/onboarding/tickets', method: 'get', params })
}

export const getTicketDetail = (ticketId) => {
  return service({ url: `/admin-api/v1/onboarding/tickets/${ticketId}`, method: 'get' })
}

export const uploadAttachment = (ticketId, data) => {
  return service({ url: `/admin-api/v1/onboarding/tickets/${ticketId}/attachments`, method: 'post', data })
}

export const submitTicket = (ticketId) => {
  return service({ url: `/admin-api/v1/onboarding/tickets/${ticketId}/submit`, method: 'put' })
}

export const reviewTicket = (ticketId, data) => {
  return service({ url: `/admin-api/v1/onboarding/tickets/${ticketId}/review`, method: 'put', data })
}
