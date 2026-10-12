import { createFileRoute } from '@tanstack/react-router'
import { Activity } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'

export const Route = createFileRoute('/_authenticated/')({
  component: DashboardPage,
})

function DashboardPage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Monitors</h1>
        <p className="text-sm text-muted-foreground">Track uptime and response times for your sites.</p>
      </div>
      <Card>
        <CardContent className="flex flex-col items-center gap-2 py-12 text-center">
          <Activity className="size-8 text-muted-foreground" />
          <p className="font-medium">No monitors yet</p>
        </CardContent>
      </Card>
    </div>
  )
}
