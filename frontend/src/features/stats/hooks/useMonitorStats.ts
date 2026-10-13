import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { useStatsRangeStore } from '@/state/useStatsRangeStore'
import { statsQueryOptions } from '../queries'

// Keeps showing the previous range while a new one loads, instead of suspending the whole page.
export function useMonitorStats(monitorId: number) {
  const range = useStatsRangeStore((state) => state.range)
  return useQuery({ ...statsQueryOptions(monitorId, range), placeholderData: keepPreviousData })
}
