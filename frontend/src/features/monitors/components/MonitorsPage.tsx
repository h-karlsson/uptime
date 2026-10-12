import { useState } from 'react'
import { useSuspenseQuery } from '@tanstack/react-query'
import { Activity, Plus } from 'lucide-react'
import type { Monitor } from '@/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { monitorsQueryOptions } from '../queries'
import { DeleteMonitorDialog } from './DeleteMonitorDialog'
import { MonitorFormDialog } from './MonitorFormDialog'
import { MonitorRow } from './MonitorRow'

// The monitor stays set while the dialog closes, so its content doesn't change mid-animation.
type DialogState = {
  open: boolean
  monitor?: Monitor
}

export function MonitorsPage() {
  const { data: monitors } = useSuspenseQuery(monitorsQueryOptions)
  const [form, setForm] = useState<DialogState>({ open: false })
  const [deletion, setDeletion] = useState<DialogState>({ open: false })

  const openForm = (monitor?: Monitor) => setForm({ open: true, monitor })

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Monitors</h1>
          <p className="text-sm text-muted-foreground">Track uptime and response times for your sites.</p>
        </div>
        <Button onClick={() => openForm()}>
          <Plus />
          Add monitor
        </Button>
      </div>

      <Card className="py-0">
        {monitors.length === 0 ? (
          <CardContent className="flex flex-col items-center gap-2 py-12 text-center">
            <Activity className="size-8 text-muted-foreground" />
            <p className="font-medium">No monitors yet</p>
            <p className="text-sm text-muted-foreground">Add a URL to start tracking its uptime.</p>
          </CardContent>
        ) : (
          <ul className="divide-y">
            {monitors.map((monitor) => (
              <MonitorRow
                key={monitor.id}
                monitor={monitor}
                onEdit={openForm}
                onDelete={(monitor) => setDeletion({ open: true, monitor })}
              />
            ))}
          </ul>
        )}
      </Card>

      <MonitorFormDialog
        open={form.open}
        onOpenChange={(open) => setForm((state) => ({ ...state, open }))}
        monitor={form.monitor}
      />
      <DeleteMonitorDialog
        open={deletion.open}
        onOpenChange={(open) => setDeletion((state) => ({ ...state, open }))}
        monitor={deletion.monitor}
      />
    </div>
  )
}
