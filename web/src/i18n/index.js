/**
 * SkyAgent Admin — i18n bootstrap (Story 8.3)
 *
 * Uses vue-i18n v9 in Composition mode (legacy: false).
 * Locale preference is persisted in localStorage under `gva_locale`.
 * Default locale is Vietnamese (vi-VN); fallback is English (en-US).
 *
 * USAGE
 *   // In <script setup>:
 *   import { useI18n } from 'vue-i18n'
 *   const { t } = useI18n()
 *   t('admin.common.create_failed')
 *
 *   // In <template>:
 *   {{ $t('admin.common.create_failed') }}
 *
 * ADDING A NEW KEY
 *   1. Add it to both locales/vi-VN.json and locales/en-US.json.
 *   2. Reference it via $t / t — never hardcode translated strings in views.
 *   3. Run `make admin-i18n-lint` to verify no CJK leaked into source.
 */
import { createI18n } from 'vue-i18n'
import viVN from './locales/vi-VN.json'
import enUS from './locales/en-US.json'

export const SUPPORTED_LOCALES = ['vi-VN', 'en-US']
export const DEFAULT_LOCALE = 'vi-VN'
export const STORAGE_KEY = 'gva_locale'

function resolveInitialLocale() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved && SUPPORTED_LOCALES.includes(saved)) return saved
  } catch (_) {
    // localStorage may be unavailable (SSR, privacy mode) — fall through.
  }
  return DEFAULT_LOCALE
}

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: resolveInitialLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'vi-VN': viVN,
    'en-US': enUS
  },
  // In dev, only warn on missing keys that are *meant* to be i18n keys
  // (the `admin.*` namespace per CLAUDE.md). Menu/breadcrumb titles stored
  // in the DB as literal display strings (e.g. "Server Status", "Trang chủ")
  // also pass through t(), and falling back to the literal is intentional.
  // Production silences warnings entirely.
  missingWarn: import.meta.env.DEV ? /^admin\./ : false,
  fallbackWarn: import.meta.env.DEV ? /^admin\./ : false
})

/**
 * Persist + apply a new locale at runtime.
 * Call from the language switcher component.
 */
export function setLocale(locale) {
  if (!SUPPORTED_LOCALES.includes(locale)) return
  i18n.global.locale.value = locale
  try {
    localStorage.setItem(STORAGE_KEY, locale)
  } catch (_) {
    /* no-op */
  }
}

export function getLocale() {
  return i18n.global.locale.value
}

export default i18n
