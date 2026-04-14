import { createSysError } from '@/api/system/sysError'

function sendErrorTip(errorInfo) {
  setTimeout(() => {
    const errorData = {
      form: errorInfo.type,
      info: `${errorInfo.message}\nStack: ${errorInfo.stack}${errorInfo.component ? `\nComponent: ${errorInfo.component.name || 'Unknown'}` : ''}${errorInfo.vueInfo ? `\nVue Info: ${errorInfo.vueInfo}` : ''}${errorInfo.source ? `\nSource: ${errorInfo.source}:${errorInfo.lineno}:${errorInfo.colno}` : ''}`,
      level: 'error',
      solution: null
    }
    
    createSysError(errorData).catch(apiErr => {
      console.error('Failed to create error record:', apiErr)
    })
  }, 0)
}
  
  window.addEventListener('unhandledrejection', (event) => {
    sendErrorTip({
      type: 'frontend',
      message: `Error: ${event.reason}`,
      stack: `Stack: ${event.reason?.stack || 'No stack trace'}`,
    });
  });
