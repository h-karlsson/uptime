import { http, HttpResponse } from 'msw'

type Credentials = { email: string; password: string }

const user = (email: string) => ({ id: 1, email, createdAt: '2026-10-12T10:00:00Z' })

// A tiny fake backend: only demo@uptime.io / demo1234 is registered.
export const authHandlers = [
  http.post('/api/auth/check-email', async ({ request }) => {
    const { email } = (await request.json()) as { email: string }
    return HttpResponse.json({ registered: email === 'demo@uptime.io' })
  }),
  http.post('/api/auth/login', async ({ request }) => {
    const { email, password } = (await request.json()) as Credentials
    return password === 'demo1234'
      ? HttpResponse.json(user(email))
      : HttpResponse.json({ error: 'invalid email or password' }, { status: 401 })
  }),
  http.post('/api/auth/register', async ({ request }) => {
    const { email } = (await request.json()) as Credentials
    return HttpResponse.json(user(email), { status: 201 })
  }),
]
