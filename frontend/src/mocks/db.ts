import type { User } from '@/client'
import { seedIncidents, seedMonitors } from './seed'

type MockUser = User & { password: string }

// In-memory data, reset on every page load. The session survives reloads so you stay logged in.
const startedAt = Date.now()

export const db = {
  users: [{ id: 1, email: 'demo@uptime.io', password: 'demo1234', createdAt: new Date(startedAt).toISOString() }] as MockUser[],
  monitors: seedMonitors().map((monitor) => ({ ...monitor, userId: 1 })),
  incidents: seedIncidents(),
  nextId: 100,
}

const SESSION_KEY = 'uptime-mock-session'

export function currentUser(): MockUser | undefined {
  const id = Number(localStorage.getItem(SESSION_KEY))
  return db.users.find((user) => user.id === id)
}

export function startSession(user: MockUser) {
  localStorage.setItem(SESSION_KEY, String(user.id))
}

export function endSession() {
  localStorage.removeItem(SESSION_KEY)
}

// Stands in for the checker: monitors added in this session come up a few seconds after creation.
export function simulateFirstChecks() {
  for (const monitor of db.monitors) {
    const createdAt = Date.parse(monitor.createdAt)
    if (monitor.status === 'unknown' && monitor.enabled && createdAt > startedAt && Date.now() - createdAt > 3000) {
      monitor.status = 'up'
    }
  }
}

export function publicUser({ id, email, createdAt }: MockUser): User {
  return { id, email, createdAt }
}
