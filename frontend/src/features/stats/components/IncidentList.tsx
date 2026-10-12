import { CircleAlert, CircleCheck } from 'lucide-react'
import type { Incident } from '@/client'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { formatDateTime, formatDuration } from '../format'

export function IncidentList({ incidents }: { incidents: Incident[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Incidents</CardTitle>
        <CardDescription>Opened after 3 failed checks in a row, resolved on the first success.</CardDescription>
      </CardHeader>
      <CardContent>
        {incidents.length === 0 ? (
          <p className="text-sm text-muted-foreground">No incidents so far.</p>
        ) : (
          <ul className="divide-y">
            {incidents.map((incident) => (
              <IncidentRow key={incident.id} incident={incident} />
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  )
}

function IncidentRow({ incident }: { incident: Incident }) {
  const { startedAt, resolvedAt, cause } = incident
  const Icon = resolvedAt === null ? CircleAlert : CircleCheck

  return (
    <li className="flex items-start gap-3 py-3 first:pt-0 last:pb-0">
      <Icon className={resolvedAt === null ? 'mt-0.5 size-4 text-red-500' : 'mt-0.5 size-4 text-muted-foreground'} />
      <div className="min-w-0 flex-1">
        <p className="text-sm font-medium">{formatDateTime(startedAt)}</p>
        {cause && <p className="truncate font-mono text-xs text-muted-foreground">{cause}</p>}
      </div>
      {resolvedAt === null ? (
        <Badge variant="destructive">Ongoing</Badge>
      ) : (
        <span className="text-sm text-muted-foreground">
          Lasted {formatDuration(new Date(resolvedAt).getTime() - new Date(startedAt).getTime())}
        </span>
      )}
    </li>
  )
}
