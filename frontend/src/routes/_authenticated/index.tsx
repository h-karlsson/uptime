import { createFileRoute } from '@tanstack/react-router'
import { MonitorsPage } from '@/features/monitors/components/MonitorsPage'
import { monitorsQueryOptions } from '@/features/monitors/queries'

export const Route = createFileRoute('/_authenticated/')({
  loader: ({ context }) => context.queryClient.ensureQueryData(monitorsQueryOptions),
  component: MonitorsPage,
})
