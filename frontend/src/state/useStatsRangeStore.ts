import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { MonitorStats } from '@/client'

export type StatsRange = MonitorStats['range']

export const STATS_RANGES: StatsRange[] = ['24h', '7d', '30d']

type StatsRangeState = {
  range: StatsRange
  setRange: (range: StatsRange) => void
}

export const useStatsRangeStore = create<StatsRangeState>()(
  persist(
    (set) => ({
      range: '24h',
      setRange: (range) => set({ range }),
    }),
    { name: 'uptime-stats-range' },
  ),
)
