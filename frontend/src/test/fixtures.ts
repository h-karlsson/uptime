import type { Monitor } from '@/client'

export function aMonitor(overrides: Partial<Monitor> = {}): Monitor {
  return {
    id: 1,
    name: 'Example',
    url: 'https://example.com/',
    method: 'GET',
    intervalSeconds: 60,
    timeoutMs: 10000,
    expectedStatus: 200,
    enabled: true,
    status: 'up',
    createdAt: '2026-10-12T10:00:00Z',
    ...overrides,
  }
}
