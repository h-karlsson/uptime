import { http, HttpResponse } from 'msw'
import { STATS_RANGES, type StatsRange } from '@/state/useStatsRangeStore'
import { db } from '../db'
import { fail, findOwnMonitor, monitorNotFound } from '../respond'
import { dailyUptime, latencySeries, summarize } from '../series'

const incidentsOf = (monitorId: number) => db.incidents.get(monitorId) ?? []

export const statsHandlers = [
  http.get('/api/monitors/:id/stats', ({ params, request }) => {
    const monitor = findOwnMonitor(params.id)
    if (!monitor) return monitorNotFound()
    const range = (new URL(request.url).searchParams.get('range') ?? '24h') as StatsRange
    if (!STATS_RANGES.includes(range)) return fail(400, 'range must be 24h, 7d or 30d')

    const latency = latencySeries(monitor, incidentsOf(monitor.id), range)
    return HttpResponse.json(summarize(monitor, latency, range))
  }),

  http.get('/api/monitors/:id/days', ({ params }) => {
    const monitor = findOwnMonitor(params.id)
    return monitor ? HttpResponse.json(dailyUptime(monitor, incidentsOf(monitor.id))) : monitorNotFound()
  }),

  http.get('/api/monitors/:id/incidents', ({ params }) => {
    const monitor = findOwnMonitor(params.id)
    return monitor ? HttpResponse.json(incidentsOf(monitor.id)) : monitorNotFound()
  }),
]
