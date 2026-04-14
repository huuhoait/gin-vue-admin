export default {
  install: (app) => {
    app.directive('click-outside', {
      mounted(el, binding) {
        const handler = (e) => {
          // If element contains target (or is removed), do not trigger
          if (!el || el.contains(e.target) || e.target === el) return
          // Supports function or object: { handler: fn, exclude: [el1, el2], capture: true }
          const value = binding.value
          if (value && typeof value === 'object') {
            if (
              value.exclude &&
              value.exclude.some(
                (ex) => ex && ex.contains && ex.contains(e.target)
              )
            )
              return
            if (typeof value.handler === 'function') value.handler(e)
          } else if (typeof value === 'function') {
            value(e)
          }
        }

        // Store on element for cleanup
        el.__clickOutsideHandler__ = handler

        // Delay registration to avoid firing during mounted click
        setTimeout(() => {
          document.addEventListener('mousedown', handler)
          document.addEventListener('touchstart', handler)
        }, 0)
      },
      unmounted(el) {
        const h = el.__clickOutsideHandler__
        if (h) {
          document.removeEventListener('mousedown', h)
          document.removeEventListener('touchstart', h)
          delete el.__clickOutsideHandler__
        }
      }
    })
  }
}
