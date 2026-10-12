import type { MonitorStats } from '@/client'
import { Card, CardContent } from '@/components/ui/card'
import { formatMs, formatPercent, formatRelative } from '../format'

export function StatCards({ stats }: { stats: MonitorStats }) {
  const cards = [
    { label: 'Uptime', value: formatPercent(stats.uptimePercent), detail: `${stats.checks} checks` },
    { label: 'Avg response', value: formatMs(stats.avgLatencyMs), detail: 'Successful checks' },
    { label: 'P95 response', value: formatMs(stats.p95LatencyMs), detail: '95% were faster' },
    { label: 'Last checked', value: formatRelative(stats.lastCheckedAt), detail: 'Refreshes automatically' },
  ]

  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
      {cards.map(({ label, value, detail }) => (
        <Card key={label} size="sm">
          <CardContent className="flex flex-col gap-1">
            <span className="text-sm text-muted-foreground">{label}</span>
            <span className="text-2xl font-semibold tracking-tight tabular-nums">{value}</span>
            <span className="text-xs text-muted-foreground">{detail}</span>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
