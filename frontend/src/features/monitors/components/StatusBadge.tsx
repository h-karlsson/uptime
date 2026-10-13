import type { Monitor } from '@/client'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

const HEARTBEAT = 'M0 6h13l1.5-2 1.5 2h2l2-5 2.5 10 2-5h3l1.5-1.5 1.5 1.5H40'
const FLATLINE = 'M0 6h40'

const STATUSES = {
  up: { label: 'Up', trace: HEARTBEAT, color: 'text-emerald-500', sweep: 'motion-safe:animate-heartbeat' },
  down: { label: 'Down', trace: FLATLINE, color: 'text-red-500', sweep: 'motion-safe:animate-flatline' },
  unknown: { label: 'Pending', trace: FLATLINE, color: 'text-muted-foreground motion-safe:animate-pulse', sweep: null },
  paused: { label: 'Paused', trace: FLATLINE, color: 'text-amber-500', sweep: null },
}

export function StatusBadge({ monitor }: { monitor: Monitor }) {
  const { label, trace, color, sweep } = STATUSES[monitor.enabled ? monitor.status : 'paused']

  return (
    <Badge variant="outline" className="min-w-28 justify-start gap-1.5 flex">
      <span className={cn('flex', color)} aria-hidden>
        <svg
          viewBox="0 0 40 12"
          className="h-3 w-10"
          fill="none"
          stroke="currentColor"
          strokeWidth={1.5}
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          {sweep && <path d={trace} className="opacity-25" />}
          <path d={trace} pathLength={100} strokeDasharray="100 100" className={sweep ?? undefined} />
        </svg>
      </span>
      <div className={"flex-1 text-center"}>
          {label}
      </div>
    </Badge>
  )
}
