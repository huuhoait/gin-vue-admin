import { fmtTitle } from '@/utils/fmtRouterTitle'
import config from '@/core/config'
import i18n from '@/i18n'
export default function getPageTitle(pageTitle, route) {
  if (pageTitle) {
    // Resolve i18n keys like "admin.router.xxx" before formatting placeholders.
    const resolved = i18n.global.te(pageTitle) ? i18n.global.t(pageTitle) : pageTitle
    const title = fmtTitle(resolved, route)
    return `${title} - ${config.appName}`
  }
  return `${config.appName}`
}
