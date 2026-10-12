import type { FormEvent } from 'react'
import type { Monitor, MonitorInput } from '@/client'
import { SubmitButton } from '@/components/SubmitButton'
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { useCreateMonitor } from '../hooks/useCreateMonitor'
import { useUpdateMonitor } from '../hooks/useUpdateMonitor'
import { DEFAULT_INPUT, INTERVAL_OPTIONS, formatInterval, toInput } from '../settings'

type MonitorFormProps = {
  monitor?: Monitor
  onDone: () => void
}

export function MonitorForm({ monitor, onDone }: MonitorFormProps) {
  const create = useCreateMonitor()
  const update = useUpdateMonitor()
  const mutation = monitor ? update : create
  const initial = monitor ? toInput(monitor) : DEFAULT_INPUT
  const intervals = INTERVAL_OPTIONS.includes(initial.intervalSeconds)
    ? INTERVAL_OPTIONS
    : [...INTERVAL_OPTIONS, initial.intervalSeconds].sort((a, b) => a - b)

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const input = readInput(new FormData(event.currentTarget))
    if (monitor) {
      update.mutate({ id: monitor.id, input }, { onSuccess: onDone })
    } else {
      create.mutate(input, { onSuccess: onDone })
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <FieldGroup>
        <Field>
          <FieldLabel htmlFor="name">Name</FieldLabel>
          <Input id="name" name="name" defaultValue={initial.name} placeholder="My website" maxLength={100} required />
        </Field>
        <Field>
          <FieldLabel htmlFor="url">URL</FieldLabel>
          <Input id="url" name="url" type="url" defaultValue={initial.url} placeholder="https://example.com" required />
        </Field>
        <div className="grid grid-cols-2 gap-4">
          <Field>
            <FieldLabel htmlFor="method">Method</FieldLabel>
            <Select name="method" defaultValue={initial.method}>
              <SelectTrigger id="method" className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="GET">GET</SelectItem>
                <SelectItem value="HEAD">HEAD</SelectItem>
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel htmlFor="intervalSeconds">Check every</FieldLabel>
            <Select name="intervalSeconds" defaultValue={String(initial.intervalSeconds)}>
              <SelectTrigger id="intervalSeconds" className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {intervals.map((seconds) => (
                  <SelectItem key={seconds} value={String(seconds)}>
                    {formatInterval(seconds)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel htmlFor="timeoutSeconds">Timeout (seconds)</FieldLabel>
            <Input
              id="timeoutSeconds"
              name="timeoutSeconds"
              type="number"
              min={1}
              max={30}
              defaultValue={initial.timeoutMs / 1000}
              required
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="expectedStatus">Expected status</FieldLabel>
            <Input
              id="expectedStatus"
              name="expectedStatus"
              type="number"
              min={100}
              max={599}
              defaultValue={initial.expectedStatus}
              required
            />
          </Field>
        </div>
        <Field orientation="horizontal">
          <Switch id="enabled" name="enabled" defaultChecked={initial.enabled} />
          <div>
            <FieldLabel htmlFor="enabled">Active</FieldLabel>
            <FieldDescription>Paused monitors aren't checked.</FieldDescription>
          </div>
        </Field>
        <FieldError>{mutation.error?.message}</FieldError>
        <SubmitButton pending={mutation.isPending}>{monitor ? 'Save changes' : 'Add monitor'}</SubmitButton>
      </FieldGroup>
    </form>
  )
}

function readInput(form: FormData): MonitorInput {
  return {
    name: String(form.get('name')),
    url: String(form.get('url')),
    method: form.get('method') === 'HEAD' ? 'HEAD' : 'GET',
    intervalSeconds: Number(form.get('intervalSeconds')),
    timeoutMs: Number(form.get('timeoutSeconds')) * 1000,
    expectedStatus: Number(form.get('expectedStatus')),
    enabled: form.get('enabled') === 'on',
  }
}
