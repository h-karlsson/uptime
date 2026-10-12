import { Link, createFileRoute, notFound } from '@tanstack/react-router'
import { MonitorDetailPage } from '@/features/monitors/components/MonitorDetailPage'
import { monitorsQueryOptions } from '@/features/monitors/queries'
import { daysQueryOptions, incidentsQueryOptions, statsQueryOptions } from '@/features/stats/queries'
import { useStatsRangeStore } from '@/state/useStatsRangeStore'

export const Route = createFileRoute('/_authenticated/monitors/$monitorId')({
  params: {
    parse: ({ monitorId }) => ({ monitorId: Number(monitorId) }),
    stringify: ({ monitorId }) => ({ monitorId: String(monitorId) }),
  },
  loader: async ({ context: { queryClient }, params: { monitorId } }) => {
    const monitors = await queryClient.ensureQueryData(monitorsQueryOptions)
    if (!monitors.some((monitor) => monitor.id === monitorId)) {
      throw notFound()
    }
    const { range } = useStatsRangeStore.getState()
    await Promise.all([
      queryClient.ensureQueryData(statsQueryOptions(monitorId, range)),
      queryClient.ensureQueryData(daysQueryOptions(monitorId)),
      queryClient.ensureQueryData(incidentsQueryOptions(monitorId)),
    ])
  },
  component: MonitorRoute,
  notFoundComponent: MonitorNotFound,
})

function MonitorRoute() {
  const { monitorId } = Route.useParams()
  return <MonitorDetailPage monitorId={monitorId} />
}

function MonitorNotFound() {
  return (
    <div className="flex flex-col items-start gap-2">
      <h1 className="text-2xl font-semibold tracking-tight">Monitor not found</h1>
      <Link to="/" className="text-sm text-muted-foreground hover:text-foreground">
        Back to monitors
      </Link>
    </div>
  )
}
