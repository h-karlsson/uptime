import type { Monitor } from '@/client'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

const STATUSES = {
  up: { label: 'Up', dot: 'bg-emerald-500' },
  down: { label: 'Down', dot: 'bg-red-500' },
  unknown: { label: 'Pending', dot: 'bg-muted-foreground/40' },
  paused: { label: 'Paused', dot: 'bg-amber-500' },
}

export function StatusBadge({ monitor }: { monitor: Monitor }) {
  const { label, dot } = STATUSES[monitor.enabled ? monitor.status : 'paused']

  return (
    <Badge variant="outline" className="gap-1.5">
      <span className={cn('size-2 rounded-full', dot)} />
      {label}
    </Badge>
  )
}
