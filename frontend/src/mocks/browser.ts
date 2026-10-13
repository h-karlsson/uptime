import { setupWorker } from 'msw/browser'
import { handlers } from './handlers'

const options = { onUnhandledFrame: 'bypass' } as const

export async function startMockWorker() {
  const worker = setupWorker(...handlers)
  await worker.start(options)

  // Browsers stop idle service workers, and a restarted worker has forgotten this tab, so it lets
  // requests through to the dev server. Re-enabling when the tab is shown again restores the mocks.
  document.addEventListener('visibilitychange', async () => {
    if (document.visibilityState !== 'visible') return
    await worker.stop()
    await worker.start(options)
  })
}
