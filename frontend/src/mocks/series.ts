import type { Incident, LatencyPoint, Monitor, MonitorStats, UptimeDay } from '@/client'
import type { StatsRange } from '@/state/useStatsRangeStore'

const MINUTE = 60_000
const DAY = 24 * 60 * MINUTE

const RANGES: Record<StatsRange, { duration: number; bucket: number }> = {
  '24h': { duration: DAY, bucket: 15 * MINUTE },
  '7d': { duration: 7 * DAY, bucket: 120 * MINUTE },
  '30d': { duration: 30 * DAY, bucket: 360 * MINUTE },
}

// Deterministic noise in [0, 1), so the same monitor always gets the same chart.
function noise(a: number, b: number): number {
  const value = Math.sin(a * 12.9898 + b * 78.233) * 43758.5453
  return value - Math.floor(value)
}

function downtimeMs(incidents: Incident[], start: number, end: number): number {
  return incidents.reduce((total, { startedAt, resolvedAt }) => {
    const from = Math.max(start, Date.parse(startedAt))
    const to = Math.min(end, resolvedAt ? Date.parse(resolvedAt) : Date.now())
    return total + Math.max(0, to - from)
  }, 0)
}

function hasChecks(monitor: Monitor): boolean {
  return monitor.status !== 'unknown'
}

export function latencySeries(monitor: Monitor, incidents: Incident[], range: StatsRange): LatencyPoint[] {
  if (!hasChecks(monitor)) return []
  const { duration, bucket } = RANGES[range]
  const first = Math.max(Date.now() - duration, Date.parse(monitor.createdAt))
  const points: LatencyPoint[] = []

  for (let start = Math.floor(first / bucket) * bucket; start < Date.now(); start += bucket) {
    const checks = Math.max(1, Math.round(bucket / (monitor.intervalSeconds * 1000)))
    const failedChecks = Math.round(checks * (downtimeMs(incidents, start, start + bucket) / bucket))
    const latency = 120 + monitor.id * 45 + noise(monitor.id, start / bucket) * 60
    points.push({
      time: new Date(start).toISOString(),
      checks,
      failedChecks,
      avgLatencyMs: failedChecks === checks ? null : Math.round(latency),
    })
  }
  return points
}

export function summarize(monitor: Monitor, latency: LatencyPoint[], range: StatsRange): MonitorStats {
  const checks = latency.reduce((sum, point) => sum + point.checks, 0)
  const failed = latency.reduce((sum, point) => sum + point.failedChecks, 0)
  const latencies = latency.flatMap((point) => (point.avgLatencyMs === null ? [] : [point.avgLatencyMs])).sort((a, b) => a - b)
  const average = latencies.reduce((sum, value) => sum + value, 0) / latencies.length

  return {
    range,
    checks,
    uptimePercent: checks ? ((checks - failed) / checks) * 100 : null,
    avgLatencyMs: latencies.length ? Math.round(average) : null,
    p95LatencyMs: latencies.length ? Math.round(latencies[Math.floor(latencies.length * 0.95)] * 1.15) : null,
    lastCheckedAt: hasChecks(monitor) ? new Date(Date.now() - 20_000).toISOString() : null,
    latency,
  }
}

export function dailyUptime(monitor: Monitor, incidents: Incident[]): UptimeDay[] {
  if (!hasChecks(monitor)) return []
  const today = Math.floor(Date.now() / DAY) * DAY
  const days: UptimeDay[] = []

  for (let start = today - 29 * DAY; start <= today; start += DAY) {
    if (start + DAY < Date.parse(monitor.createdAt)) continue
    const end = Math.min(start + DAY, Date.now())
    days.push({
      date: new Date(start).toISOString().slice(0, 10),
      checks: Math.round((end - start) / (monitor.intervalSeconds * 1000)),
      uptimePercent: 100 - (downtimeMs(incidents, start, end) / (end - start)) * 100,
    })
  }
  return days
}
