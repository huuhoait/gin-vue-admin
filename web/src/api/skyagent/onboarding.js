import service from '@/utils/request'

// @Summary Create Agent L0 with full info + auto ticket
// @Router /admin-api/v1/onboarding/agents [post]
export const onboardingAgent = (data) => {
  return service({
    url: '/admin-api/v1/onboarding/agents',
    method: 'post',
    data
  })
}
