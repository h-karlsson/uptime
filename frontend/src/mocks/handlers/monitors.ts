import { http, HttpResponse } from 'msw'
import type { MonitorInput } from '@/client'
import { currentUser, db, simulateFirstChecks } from '../db'
import { fail, findOwnMonitor, monitorNotFound, notLoggedIn, toMonitor } from '../respond'

const MAX_MONITORS = 10

// The main rules of the real backend: an http(s) url, a name, and one monitor per url.
function rejectInput(input: MonitorInput, userId: number, exceptId?: number) {
  if (!input.name.trim()) return fail(400, 'name is required')
  if (!/^https?:\/\/[^/]+/.test(input.url)) return fail(400, 'url must be an http or https address')
  const duplicate = db.monitors.some(
    (monitor) => monitor.userId === userId && monitor.id !== exceptId && monitor.url === input.url,
  )
  return duplicate ? fail(409, 'you already have a monitor for this url') : null
}

export const monitorHandlers = [
  http.get('/api/monitors', () => {
    const user = currentUser()
    if (!user) return notLoggedIn()
    simulateFirstChecks()
    return HttpResponse.json(db.monitors.filter((monitor) => monitor.userId === user.id).map(toMonitor))
  }),

  http.post('/api/monitors', async ({ request }) => {
    const user = currentUser()
    if (!user) return notLoggedIn()
    if (db.monitors.filter((monitor) => monitor.userId === user.id).length >= MAX_MONITORS) {
      return fail(409, `you can have at most ${MAX_MONITORS} monitors`)
    }
    const input = (await request.json()) as MonitorInput
    const rejection = rejectInput(input, user.id)
    if (rejection) return rejection

    const monitor = { ...input, id: db.nextId++, userId: user.id, status: 'unknown' as const, createdAt: new Date().toISOString() }
    db.monitors.push(monitor)
    return HttpResponse.json(toMonitor(monitor), { status: 201 })
  }),

  http.get('/api/monitors/:id', ({ params }) => {
    const monitor = findOwnMonitor(params.id)
    return monitor ? HttpResponse.json(toMonitor(monitor)) : monitorNotFound()
  }),

  http.put('/api/monitors/:id', async ({ params, request }) => {
    const monitor = findOwnMonitor(params.id)
    if (!monitor) return monitorNotFound()
    const input = (await request.json()) as MonitorInput
    const rejection = rejectInput(input, monitor.userId, monitor.id)
    if (rejection) return rejection

    Object.assign(monitor, input)
    return HttpResponse.json(toMonitor(monitor))
  }),

  http.delete('/api/monitors/:id', ({ params }) => {
    const monitor = findOwnMonitor(params.id)
    if (!monitor) return monitorNotFound()
    db.monitors.splice(db.monitors.indexOf(monitor), 1)
    return new HttpResponse(null, { status: 204 })
  }),
]
