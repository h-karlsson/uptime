import { Ellipsis } from 'lucide-react'
import type { Monitor } from '@/client'
import { Card } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { StatusBadge } from '@/features/monitors/components/StatusBadge'
import { formatInterval } from '@/features/monitors/settings'

function sample(id: number, name: string, url: string, overrides: Partial<Monitor> = {}): Monitor {
  return {
    id,
    name,
    url,
    method: 'GET',
    intervalSeconds: 60,
    timeoutMs: 10000,
    expectedStatus: 200,
    enabled: true,
    status: 'up',
    createdAt: '2026-01-01T00:00:00Z',
    ...overrides,
  }
}

const MONITORS = [
  sample(1, 'Website', 'https://example.com/'),
  sample(2, 'API', 'https://api.example.com/health', { intervalSeconds: 30 }),
  sample(3, 'Docs', 'https://docs.example.com/', { status: 'down' }),
  sample(4, 'Staging', 'https://staging.example.com/', { enabled: false }),
  sample(5, 'Blog', 'https://blog.example.com/', { status: 'unknown' }),
]

// A static copy of the monitors page, so visitors see the product before signing in.
export function MonitorPreview() {
  return (
    <div inert className="flex max-w-xl flex-col gap-4">
      <Card className="py-0">
        <ul className="divide-y">
          {MONITORS.map((monitor) => (
            <li key={monitor.id} className="flex items-center gap-4 px-4 py-3">
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium">{monitor.name}</p>
                <p className="truncate text-sm text-muted-foreground">{monitor.url}</p>
              </div>
              <span className="text-sm text-muted-foreground">Every {formatInterval(monitor.intervalSeconds)}</span>
              <StatusBadge monitor={monitor} />
              <Switch defaultChecked={monitor.enabled} />
              <Ellipsis className="mx-1.5 size-4" />
            </li>
          ))}
        </ul>
      </Card>
    </div>
  )
}
