import type { UptimeDay } from '@/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'
import { formatDate, formatPercent } from '../format'

const DAYS = 30

// One slot per UTC day, oldest first, matching the backend's daily history.
function lastDays(): string[] {
  const today = new Date()
  return Array.from({ length: DAYS }, (_, index) => {
    const day = new Date(Date.UTC(today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate() - (DAYS - 1 - index)))
    return day.toISOString().slice(0, 10)
  })
}

function barColor(day: UptimeDay | undefined): string {
  if (!day) return 'bg-muted'
  if (day.uptimePercent >= 99) return 'bg-emerald-500'
  if (day.uptimePercent >= 95) return 'bg-amber-500'
  return 'bg-red-500'
}

export function UptimeBars({ days }: { days: UptimeDay[] }) {
  const byDate = new Map(days.map((day) => [day.date, day]))
  const checks = days.reduce((sum, day) => sum + day.checks, 0)
  const okChecks = days.reduce((sum, day) => sum + (day.checks * day.uptimePercent) / 100, 0)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Last 30 days</CardTitle>
        <CardDescription>{checks > 0 ? `${formatPercent((okChecks / checks) * 100)} uptime` : 'No checks yet'}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-2">
        <div className="flex h-10 gap-0.5 sm:gap-1">
          {lastDays().map((date) => {
            const day = byDate.get(date)
            return (
              <Tooltip key={date}>
                <TooltipTrigger asChild>
                  <div
                    className={cn('flex-1 rounded-sm transition-opacity hover:opacity-70', barColor(day))}
                    aria-label={`${formatDate(date)}: ${day ? formatPercent(day.uptimePercent) : 'no data'}`}
                  />
                </TooltipTrigger>
                <TooltipContent>
                  <p className="font-medium">{formatDate(date)}</p>
                  <p>{day ? `${formatPercent(day.uptimePercent)} uptime · ${day.checks} checks` : 'No data'}</p>
                </TooltipContent>
              </Tooltip>
            )
          })}
        </div>
        <div className="flex justify-between text-xs text-muted-foreground">
          <span>30 days ago</span>
          <span>Today</span>
        </div>
      </CardContent>
    </Card>
  )
}
