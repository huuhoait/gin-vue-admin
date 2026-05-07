import service from '@/utils/request'

// Drop empty optional fields so the payload matches the omitempty semantics
// in OnboardingAgentRequest (external-frontend-integration.md §14.2, §14.6
// rule 1 — never re-submit masked PII).
const pruneEmpty = (obj) => {
  const out = {}
  for (const [k, v] of Object.entries(obj)) {
    if (v !== '' && v !== null && v !== undefined) out[k] = v
  }
  return out
}

// @Summary Create Agent L0 with full info + auto ticket
// @Router /admin-api/v1/onboarding/agents [post]
// X-Idempotency-Key is auto-injected by utils/request.js for any
// /admin-api/v1/* mutation (§6). Pass `idempotencyKey` only to pin one
// key across a user-driven retry sequence.
export const onboardingAgent = (data, idempotencyKey) => {
  return service({
    url: '/admin-api/v1/onboarding/agents',
    method: 'post',
    data: pruneEmpty(data),
    ...(idempotencyKey ? { headers: { 'X-Idempotency-Key': idempotencyKey } } : {})
  })
}
