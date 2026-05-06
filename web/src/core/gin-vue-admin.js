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
      
    `)
  }
}
