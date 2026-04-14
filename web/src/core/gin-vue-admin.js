/*
 * gin-vue-admin web framework
 */
// Load website config and register plugins
import { register } from './global'
import packageInfo from '../../package.json'

export default {
  install: (app) => {
    register(app)
    console.log(`
       Welcome to Gin-Vue-Admin
       Version: v${packageInfo.version}
       Community: WeChat shouzi_1994 | QQ group 622360840
       Repo: https://github.com/flipped-aurora/gin-vue-admin
       Plugin market: https://plugin.gin-vue-admin.com
       Community forum: https://support.qq.com/products/371961
       Swagger: http://127.0.0.1:${import.meta.env.VITE_SERVER_PORT}/swagger/index.html
       Dev server: http://127.0.0.1:${import.meta.env.VITE_CLI_PORT}
       Support the project: https://www.gin-vue-admin.com/coffee/index.html
       -------------------------------------- LICENSE NOTICE --------------------------------------
       ** Owner: flipped-aurora OSS team **
       ** Company: Flipped Aurora (Beijing) Technology Co., Ltd. **
       ** Commercial use requires a license: https://plugin.gin-vue-admin.com/license **
       ** Thanks for supporting Gin-Vue-Admin. Licensed usage helps long-term maintenance. **
    `)
  }
}
