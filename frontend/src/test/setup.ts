import '@testing-library/jest-dom/vitest'
import { cleanup } from '@testing-library/react'
import { afterAll, afterEach, beforeAll } from 'vitest'
import { client } from '@/client/client.gen'
import { setupApiClient } from '@/lib/api'
import { server } from './server'

client.setConfig({ baseUrl: window.location.origin })
setupApiClient()

beforeAll(() => server.listen({ onUnhandledFrame: 'error' }))

afterEach(() => {
  server.resetHandlers()
  cleanup()
  localStorage.clear()
})

afterAll(() => server.close())
