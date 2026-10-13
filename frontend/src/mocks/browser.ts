import { setupWorker } from 'msw/browser'
import { handlers } from './handlers'

export async function startMockWorker() {
  await setupWorker(...handlers).start({ onUnhandledFrame: 'bypass' })

  window.addEventListener('focus', reloadIfMocksStopped)
  document.addEventListener('visibilitychange', reloadIfMocksStopped)
}

async function reloadIfMocksStopped() {
  if (document.visibilityState !== 'visible') return
  const response = await fetch('/api/health')
  if (!response.ok) location.reload()
}
