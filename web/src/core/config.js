/**
 * Site config (Story 8.3: Chinese banner converted to English)
 */
import packageInfo from '../../package.json'

const greenText = (text) => `\x1b[32m${text}\x1b[0m`

export const config = {
  appName: 'SkyAgent Admin',
  showViteLogo: true,
  keepAliveTabs: false,
  logs: []
}

// Story 8.3: CLI banner converted from Chinese to English. UI-visible text is
// handled by vue-i18n (see src/i18n/); logs remain English per the monorepo
// logging convention.
export const viteLogo = (env) => {
  if (config.showViteLogo) {
    console.log(greenText(`> SkyAgent Admin Portal (forked from gin-vue-admin)`))
    console.log(greenText(`> Upstream: https://github.com/flipped-aurora/gin-vue-admin`))
    console.log(greenText(`> Version: v${packageInfo.version}`))
    console.log(
      greenText(
        `> Swagger: http://127.0.0.1:${env.VITE_SERVER_PORT}/swagger/index.html`
      )
    )
    console.log(
      greenText(`> Dev server: http://127.0.0.1:${env.VITE_CLI_PORT}`)
    )
    console.log('\n')
  }
}

export default config
