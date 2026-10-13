import { useState } from 'react'
import { useSuspenseQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { ArrowLeft, Pencil } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { IncidentList } from '@/features/stats/components/IncidentList'
import { LatencyChart } from '@/features/stats/components/LatencyChart'
import { RangePicker } from '@/features/stats/components/RangePicker'
import { StatCards } from '@/features/stats/components/StatCards'
import { UptimeBars } from '@/features/stats/components/UptimeBars'
import { useMonitorStats } from '@/features/stats/hooks/useMonitorStats'
import { daysQueryOptions, incidentsQueryOptions } from '@/features/stats/queries'
import { cn } from '@/lib/utils'
import { monitorsQueryOptions } from '../queries'
import { formatInterval } from '../settings'
import { MonitorFormDialog } from './MonitorFormDialog'
import { StatusBadge } from './StatusBadge'

export function MonitorDetailPage({ monitorId }: { monitorId: number }) {
  const { data: monitor } = useSuspenseQuery({
    ...monitorsQueryOptions,
    select: (monitors) => monitors.find((candidate) => candidate.id === monitorId),
  })
  const { data: stats, isPlaceholderData: changingRange } = useMonitorStats(monitorId)
  const { data: days } = useSuspenseQuery(daysQueryOptions(monitorId))
  const { data: incidents } = useSuspenseQuery(incidentsQueryOptions(monitorId))
  const [editing, setEditing] = useState(false)

  if (!monitor) {
    return <p className="text-muted-foreground">This monitor no longer exists.</p>
  }

  return (
    <div className="flex flex-col gap-6">
      <Link to="/" className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="size-4" />
        Monitors
      </Link>

      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="min-w-0">
          <div className="flex items-center gap-3">
            <h1 className="truncate text-2xl font-semibold tracking-tight">{monitor.name}</h1>
            <StatusBadge monitor={monitor} />
          </div>
          <p className="truncate text-sm text-muted-foreground">
            {monitor.method} {monitor.url} · every {formatInterval(monitor.intervalSeconds)}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <RangePicker />
          <Button variant="outline" size="sm" onClick={() => setEditing(true)}>
            <Pencil />
            Edit
          </Button>
        </div>
      </div>

      {stats && (
        <div className={cn('flex flex-col gap-6 transition-opacity', changingRange && 'opacity-60')}>
          <StatCards stats={stats} />
          <LatencyChart stats={stats} />
        </div>
      )}
      <UptimeBars days={days} />
      <IncidentList incidents={incidents} />

      <MonitorFormDialog open={editing} onOpenChange={setEditing} monitor={monitor} />
    </div>
  )
}
