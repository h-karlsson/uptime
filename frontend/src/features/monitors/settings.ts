import type { Monitor, MonitorInput } from '@/client'

export const INTERVAL_OPTIONS = [30, 60, 300, 600, 1800, 3600, 86400]

export const DEFAULT_INPUT: MonitorInput = {
  name: '',
  url: '',
  method: 'GET',
  intervalSeconds: 60,
  timeoutMs: 10000,
  expectedStatus: 200,
  enabled: true,
}

export function toInput({ name, url, method, intervalSeconds, timeoutMs, expectedStatus, enabled }: Monitor): MonitorInput {
  return { name, url, method, intervalSeconds, timeoutMs, expectedStatus, enabled }
}

export function formatInterval(seconds: number): string {
  if (seconds % 86400 === 0) return `${seconds / 86400}d`
  if (seconds % 3600 === 0) return `${seconds / 3600}h`
  if (seconds % 60 === 0) return `${seconds / 60}m`
  return `${seconds}s`
}
