import config from './config'
import { h } from 'vue'

// Import Element Plus icon components
import * as ElIconModules from '@element-plus/icons-vue'
import svgIcon from '@/components/svgIcon/svgIcon.vue'
// Icon name conversion helpers (if needed)

const createIconComponent = (name) => ({
  name: 'SvgIcon',
  render() {
    return h(svgIcon, {
      localIcon: name
    })
  }
})

const registerIcons = async (app) => {
  const iconModules = import.meta.glob('@/assets/icons/**/*.svg') // app svg icons
  const pluginIconModules = import.meta.glob(
    '@/plugin/**/assets/icons/**/*.svg'
  ) // plugin svg icons
  const mergedIconModules = Object.assign({}, iconModules, pluginIconModules) // merged svg icons
  let allKeys = []
  for (const path in mergedIconModules) {
    let pluginName = ''
    if (path.startsWith('/src/plugin/')) {
      pluginName = `${path.split('/')[3]}-`
    }
    const iconName = path
      .split('/')
      .pop()
      .replace(/\.svg$/, '')
    // Skip invalid icon names (contains whitespace)
    if (iconName.indexOf(' ') !== -1) {
      console.error(`icon ${iconName}.svg includes whitespace in ${path}`)
      continue
    }
    const key = `${pluginName}${iconName}`
    const iconComponent = createIconComponent(key)
    config.logs.push({
      key: key,
      label: key
    })
    app.component(key, iconComponent)

    // In dev mode, list all icons for easy copy/paste
    allKeys.push(key)
  }

  import.meta.env.MODE == 'development' &&
    console.log(`All available local icons: ${allKeys.join(', ')}`)
}

export const register = (app) => {
  // Register Element Plus icons
  for (const iconName in ElIconModules) {
    app.component(iconName, ElIconModules[iconName])
  }
  app.component('SvgIcon', svgIcon)
  registerIcons(app)
  app.config.globalProperties.$GIN_VUE_ADMIN = config
}
