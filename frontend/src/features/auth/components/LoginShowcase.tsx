import { MonitorPreview } from './MonitorPreview'

const FACTS = [
  { value: '30s', label: 'Shortest check interval' },
  { value: '3', label: 'Failed checks open an incident' },
  { value: '30d', label: 'Of daily history' },
]

export function LoginShowcase() {
  return (
    <div className="hidden p-4 lg:flex">
      <section
        aria-label="About Uptime"
        className="dark flex flex-1 flex-col justify-between gap-10 rounded-3xl bg-background p-12 text-foreground dark:ring-1 dark:ring-border"
      >
        <div className="flex max-w-lg flex-col gap-4">
          <span className="flex w-fit items-center gap-2 rounded-full border px-3 py-1 text-xs text-muted-foreground">
            <span className="size-1.5 rounded-full bg-emerald-500 ring-3 ring-emerald-500/20" />
            Watching your sites around the clock
          </span>
          <h2 className="text-4xl font-semibold tracking-tight text-balance">Know the moment your site goes down.</h2>
          <p className="text-base text-muted-foreground">
            Uptime checks your URLs on a schedule, opens an incident when they fail and keeps 30 days of history for
            every monitor.
          </p>
        </div>
        <MonitorPreview />
        <ul className="grid max-w-xl grid-cols-3 gap-6">
          {FACTS.map(({ value, label }) => (
            <li key={label} className="flex flex-col">
              <span className="text-2xl font-semibold tabular-nums">{value}</span>
              <span className="text-sm text-muted-foreground">{label}</span>
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}
