import { HttpResponse } from 'msw'
import type { Monitor } from '@/client'
import { currentUser, db, simulateFirstChecks } from './db'

export function fail(status: number, error: string) {
  return HttpResponse.json({ error }, { status })
}

export const notLoggedIn = () => fail(401, 'not logged in')
export const monitorNotFound = () => fail(404, 'monitor not found')

// Returns the logged-in user's monitor with this id, like the real API's ownership check.
export function findOwnMonitor(id: unknown) {
  simulateFirstChecks()
  const user = currentUser()
  return db.monitors.find((monitor) => monitor.id === Number(id) && monitor.userId === user?.id)
}

export function toMonitor({ userId: _userId, ...monitor }: Monitor & { userId: number }): Monitor {
  return monitor
}
