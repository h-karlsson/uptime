import { Ellipsis, Pencil, Trash2 } from 'lucide-react'
import type { Monitor } from '@/client'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Switch } from '@/components/ui/switch'
import { useUpdateMonitor } from '../hooks/useUpdateMonitor'
import { formatInterval, toInput } from '../settings'
import { StatusBadge } from './StatusBadge'

type MonitorRowProps = {
  monitor: Monitor
  onEdit: (monitor: Monitor) => void
  onDelete: (monitor: Monitor) => void
}

export function MonitorRow({ monitor, onEdit, onDelete }: MonitorRowProps) {
  const update = useUpdateMonitor()

  function setEnabled(enabled: boolean) {
    update.mutate({ id: monitor.id, input: { ...toInput(monitor), enabled } })
  }

  return (
    <li className="flex items-center gap-4 px-4 py-3">
      <div className="min-w-0 flex-1">
        <p className="truncate font-medium">{monitor.name}</p>
        <a
          href={monitor.url}
          target="_blank"
          rel="noreferrer"
          className="block truncate text-sm text-muted-foreground hover:underline"
        >
          {monitor.url}
        </a>
      </div>
      <span className="hidden text-sm text-muted-foreground sm:inline">
        Every {formatInterval(monitor.intervalSeconds)}
      </span>
      <StatusBadge monitor={monitor} />
      <Switch
        checked={monitor.enabled}
        onCheckedChange={setEnabled}
        aria-label={`${monitor.enabled ? 'Pause' : 'Resume'} ${monitor.name}`}
      />
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon-sm" aria-label={`Actions for ${monitor.name}`}>
            <Ellipsis />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onSelect={() => onEdit(monitor)}>
            <Pencil />
            Edit
          </DropdownMenuItem>
          <DropdownMenuItem variant="destructive" onSelect={() => onDelete(monitor)}>
            <Trash2 />
            Delete
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </li>
  )
}
