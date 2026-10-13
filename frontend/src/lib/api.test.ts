import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import * as sdk from '@/client'
import { server } from '@/test/server'
import { ApiError } from './api'

async function errorFromHealthCheck() {
  const error = await sdk.getHealth().catch((error: unknown) => error)
  expect(error).toBeInstanceOf(ApiError)
  return error as ApiError
}

describe('API client errors', () => {
  it('uses the message and status from the error body', async () => {
    server.use(http.get('/api/health', () => HttpResponse.json({ error: 'database down' }, { status: 503 })))

    const error = await errorFromHealthCheck()
    expect(error.status).toBe(503)
    expect(error.message).toBe('database down')
  })

  it('falls back to the status when the body has no message', async () => {
    server.use(http.get('/api/health', () => new HttpResponse('Bad Gateway', { status: 502 })))

    expect((await errorFromHealthCheck()).message).toBe('Request failed (502)')
  })

  it('reports network failures with status 0', async () => {
    server.use(http.get('/api/health', () => HttpResponse.error()))

    const error = await errorFromHealthCheck()
    expect(error.status).toBe(0)
    expect(error.message).toBe('Could not reach the server')
  })
})
