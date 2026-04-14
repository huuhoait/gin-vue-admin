import './style/element_visiable.scss'
import 'element-plus/theme-chalk/dark/css-vars.css'
import 'uno.css'
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import { setupVueRootValidator } from 'vite-check-multiple-dom/client';

import 'element-plus/dist/index.css'
// gin-vue-admin frontend bootstrap
import './core/gin-vue-admin'
// Router
import router from '@/router/index'
import '@/permission'
import run from '@/core/gin-vue-admin.js'
import auth from '@/directive/auth'
import clickOutSide from '@/directive/clickOutSide'
import { store } from '@/pinia'
import App from './App.vue'
import '@/core/error-handel'
// Story 8.3: i18n bootstrap. Must be installed before components that call $t.
import i18n from '@/i18n'

const app = createApp(App)

app.config.productionTip = false

setupVueRootValidator(app, {
    lang: 'en'
  })

app
  .use(run)
  .use(ElementPlus)
  .use(store)
  .use(i18n)
  .use(auth)
  .use(clickOutSide)
  .use(router)
  .mount('#app')
export default app
