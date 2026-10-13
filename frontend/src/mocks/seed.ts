import type { Incident, Monitor } from '@/client'

const minutesAgo = (minutes: number) => new Date(Date.now() - minutes * 60_000).toISOString()
const DAY = 24 * 60

const defaults = { method: 'GET', intervalSeconds: 60, timeoutMs: 10000, expectedStatus: 200, enabled: true } as const

export function seedMonitors(): Monitor[] {
  return [
    { ...defaults, id: 1, name: 'Website', url: 'https://example.com/', status: 'up', createdAt: minutesAgo(45 * DAY) },
    { ...defaults, id: 2, name: 'API', url: 'https://api.example.com/health', intervalSeconds: 30, status: 'up', createdAt: minutesAgo(40 * DAY) },
    { ...defaults, id: 3, name: 'Docs', url: 'https://docs.example.com/', status: 'down', createdAt: minutesAgo(35 * DAY) },
    { ...defaults, id: 4, name: 'Staging', url: 'https://staging.example.com/', enabled: false, status: 'up', createdAt: minutesAgo(20 * DAY) },
    { ...defaults, id: 5, name: 'Blog', url: 'https://blog.example.com/', status: 'unknown', createdAt: minutesAgo(1) },
  ]
}

function incident(id: number, startedMinutesAgo: number, lastedMinutes: number | null, cause: string): Incident {
  return {
    id,
    startedAt: minutesAgo(startedMinutesAgo),
    resolvedAt: lastedMinutes === null ? null : minutesAgo(startedMinutesAgo - lastedMinutes),
    cause,
  }
}

// Newest first, like the real API.
export function seedIncidents(): Map<number, Incident[]> {
  return new Map([
    [1, [incident(1, 3 * DAY, 14, 'timed out after 10s')]],
    [2, [incident(2, 9 * DAY, 6, 'expected status 200, got 502'), incident(3, 21 * DAY, 95, 'timed out after 10s')]],
    [3, [incident(4, 18, null, 'expected status 200, got 503'), incident(5, 12 * DAY, 41, 'could not resolve host docs.example.com')]],
    [4, [incident(6, 6 * DAY, 180, 'connect: connection refused')]],
  ])
}
