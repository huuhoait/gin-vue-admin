import i18n from '@/i18n'

// translateError maps an envelope or error-like object onto a user-facing
// Vietnamese/English string via the `errors.<responseCode>` key catalogue
// (external-frontend-integration.md §5, §11). Falls back to the envelope
// `msg` field, then `errors.unknown`.
export const translateError = (envOrErr) => {
  const t = i18n.global.t

  const env = envOrErr?.data ? envOrErr : envOrErr?.response?.data
  const responseCode = env?.data?.responseCode ?? env?.code

  if (typeof responseCode === 'number' && responseCode !== 0) {
    const key = `errors.${responseCode}`
    const translated = t(key)
    if (translated !== key) return translated
  }

  return env?.msg || t('errors.unknown')
}

// extractFieldErrors turns a 422 envelope's data.details[] into a
// { field: message } map suitable for binding to Element Plus form errors.
export const extractFieldErrors = (envOrErr) => {
  const env = envOrErr?.data ? envOrErr : envOrErr?.response?.data
  const details = env?.data?.details
  if (!Array.isArray(details)) return {}
  return Object.fromEntries(details.map((d) => [d.field, d.message]))
}

// isRetryable returns true for envelope response codes the BFF contract
// marks as safe to auto-retry with backoff (§5.x retry cheat-sheet).
const RETRYABLE_CODES = new Set([
  1003, 2003, 3001, 5001, 5002, 5003, 5004,
  11007, 15001, 21003, 32001, 52001, 52002
])
export const isRetryable = (envOrErr) => {
  const env = envOrErr?.data ? envOrErr : envOrErr?.response?.data
  const code = env?.data?.responseCode ?? env?.code
  return RETRYABLE_CODES.has(code)
}
