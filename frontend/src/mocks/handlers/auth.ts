import { http, HttpResponse } from 'msw'
import type { Credentials } from '@/client'
import { currentUser, db, endSession, publicUser, startSession } from '../db'
import { fail, notLoggedIn } from '../respond'

const normalize = (email: string) => email.trim().toLowerCase()

export const authHandlers = [
  http.post('/api/auth/check-email', async ({ request }) => {
    const { email } = (await request.json()) as { email: string }
    return HttpResponse.json({ registered: db.users.some((user) => user.email === normalize(email)) })
  }),

  http.post('/api/auth/login', async ({ request }) => {
    const { email, password } = (await request.json()) as Credentials
    const user = db.users.find((candidate) => candidate.email === normalize(email) && candidate.password === password)
    if (!user) return fail(401, 'invalid email or password')
    startSession(user)
    return HttpResponse.json(publicUser(user))
  }),

  http.post('/api/auth/register', async ({ request }) => {
    const { email, password } = (await request.json()) as Credentials
    if (db.users.some((user) => user.email === normalize(email))) return fail(409, 'email is already registered')
    if (password.length < 8) return fail(400, 'password must be at least 8 characters')
    const user = { id: db.nextId++, email: normalize(email), password, createdAt: new Date().toISOString() }
    db.users.push(user)
    startSession(user)
    return HttpResponse.json(publicUser(user), { status: 201 })
  }),

  http.post('/api/auth/logout', () => {
    endSession()
    return new HttpResponse(null, { status: 204 })
  }),

  http.get('/api/auth/me', () => {
    const user = currentUser()
    return user ? HttpResponse.json(publicUser(user)) : notLoggedIn()
  }),
]
