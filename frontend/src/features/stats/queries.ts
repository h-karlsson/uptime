import { queryOptions } from '@tanstack/react-query'
import type { StatsRange } from '@/state/useStatsRangeStore'
import { getMonitorStats, listMonitorDays, listMonitorIncidents } from './api'

export const statsQueryOptions = (monitorId: number, range: StatsRange) =>
  queryOptions({
    queryKey: ['monitors', monitorId, 'stats', range],
    queryFn: ({ signal }) => getMonitorStats(monitorId, range, signal),
    staleTime: 10_000,
    refetchInterval: 30_000,
  })

export const daysQueryOptions = (monitorId: number) =>
  queryOptions({
    queryKey: ['monitors', monitorId, 'days'],
    queryFn: ({ signal }) => listMonitorDays(monitorId, signal),
    staleTime: 10_000,
    refetchInterval: 60_000,
  })

export const incidentsQueryOptions = (monitorId: number) =>
  queryOptions({
    queryKey: ['monitors', monitorId, 'incidents'],
    queryFn: ({ signal }) => listMonitorIncidents(monitorId, signal),
    staleTime: 10_000,
    refetchInterval: 30_000,
  })
