import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import { ApiError } from '@/lib/api'
import { server } from '@/test/server'
import { getMe } from './api'

const demoUser = { id: 1, email: 'demo@uptime.io', createdAt: '2026-10-12T10:00:00Z' }

describe('getMe', () => {
  it('returns the logged-in user', async () => {
    server.use(http.get('/api/auth/me', () => HttpResponse.json(demoUser)))

    await expect(getMe()).resolves.toEqual(demoUser)
  })

  it('returns null when nobody is logged in', async () => {
    server.use(http.get('/api/auth/me', () => HttpResponse.json({ error: 'not logged in' }, { status: 401 })))

    await expect(getMe()).resolves.toBeNull()
  })

  it('still fails on other errors', async () => {
    server.use(http.get('/api/auth/me', () => HttpResponse.json({ error: 'Internal server error' }, { status: 500 })))

    await expect(getMe()).rejects.toBeInstanceOf(ApiError)
  })
})
