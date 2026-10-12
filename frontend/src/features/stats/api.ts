import * as sdk from '@/client'
import type { StatsRange } from '@/state/useStatsRangeStore'

export async function getMonitorStats(monitorId: number, range: StatsRange, signal?: AbortSignal) {
  const { data } = await sdk.getMonitorStats({ path: { id: monitorId }, query: { range }, signal })
  return data
}

export async function listMonitorDays(monitorId: number, signal?: AbortSignal) {
  const { data } = await sdk.listMonitorDays({ path: { id: monitorId }, signal })
  return data
}

export async function listMonitorIncidents(monitorId: number, signal?: AbortSignal) {
  const { data } = await sdk.listMonitorIncidents({ path: { id: monitorId }, signal })
  return data
}
