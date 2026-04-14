// Permission-based button display directive
import { useUserStore } from '@/pinia/modules/user'
export default {
  install: (app) => {
    const userStore = useUserStore()
    app.directive('auth', {
      // When the bound element is inserted into the DOM...
      mounted: function (el, binding) {
        const userInfo = userStore.userInfo
        if (!binding.value){
          el.parentNode.removeChild(el)
          return
        }
        const waitUse = binding.value.toString().split(',')
        let flag = waitUse.some((item) => Number(item) === userInfo.authorityId)
        if (binding.modifiers.not) {
          flag = !flag
        }
        if (!flag) {
          el.parentNode.removeChild(el)
        }
      }
    })
  }
}
