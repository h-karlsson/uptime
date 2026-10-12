import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts'
import type { MonitorStats } from '@/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ChartContainer, ChartTooltip, ChartTooltipContent, type ChartConfig } from '@/components/ui/chart'
import { formatChartTick, formatDateTime } from '../format'

const chartConfig = {
  avgLatencyMs: { label: 'Response time', color: 'var(--color-emerald-500)' },
} satisfies ChartConfig

export function LatencyChart({ stats }: { stats: MonitorStats }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Response time</CardTitle>
        <CardDescription>Average of successful checks; gaps are periods where every check failed.</CardDescription>
      </CardHeader>
      <CardContent>
        {stats.latency.length === 0 ? (
          <p className="flex h-56 items-center justify-center text-sm text-muted-foreground">
            No checks in this range yet.
          </p>
        ) : (
          <ChartContainer config={chartConfig} className="h-56 w-full">
            <AreaChart data={stats.latency} margin={{ left: 0, right: 8 }}>
              <defs>
                <linearGradient id="latencyFill" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="var(--color-avgLatencyMs)" stopOpacity={0.4} />
                  <stop offset="95%" stopColor="var(--color-avgLatencyMs)" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="time"
                tickLine={false}
                axisLine={false}
                minTickGap={32}
                tickFormatter={(time: string) => formatChartTick(time, stats.range)}
              />
              <YAxis tickLine={false} axisLine={false} width={64} tickFormatter={(value: number) => `${value} ms`} />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    labelFormatter={(_, payload) => formatDateTime(payload[0]?.payload.time)}
                    formatter={(value) => (
                      <div className="flex w-full justify-between gap-4">
                        <span className="text-muted-foreground">Response time</span>
                        <span className="font-mono font-medium tabular-nums">{value} ms</span>
                      </div>
                    )}
                  />
                }
              />
              <Area
                dataKey="avgLatencyMs"
                type="monotone"
                stroke="var(--color-avgLatencyMs)"
                strokeWidth={2}
                fill="url(#latencyFill)"
                dot={stats.latency.length <= 3 && { r: 3, fill: 'var(--color-avgLatencyMs)' }}
                isAnimationActive={false}
              />
            </AreaChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  )
}
