import type { Router } from 'vue-router'

export function enableRouterDebug(router?: Router) {
  if (import.meta.env.DEV) {
    console.log('[Router Debug] Enabled')
    
    // You can add router debugging logic here
    // For example, log all route changes:
    if (router) {
      router.beforeEach((to, from) => {
        console.log('[Router]', from.path, '->', to.path)
      })
    }
  }
}
