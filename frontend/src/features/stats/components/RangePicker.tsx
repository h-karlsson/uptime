import type { TransitionStartFunction } from 'react'
import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group'
import { STATS_RANGES, useStatsRangeStore, type StatsRange } from '@/state/useStatsRangeStore'

// Switching inside a transition keeps the current stats on screen while the new range loads.
export function RangePicker({ startTransition }: { startTransition: TransitionStartFunction }) {
  const range = useStatsRangeStore((state) => state.range)
  const setRange = useStatsRangeStore((state) => state.setRange)

  return (
    <ToggleGroup
      type="single"
      variant="outline"
      size="sm"
      value={range}
      onValueChange={(value) => value && startTransition(() => setRange(value as StatsRange))}
      aria-label="Time range"
    >
      {STATS_RANGES.map((option) => (
        <ToggleGroupItem key={option} value={option} className="px-3">
          {option}
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
  )
}
