import { Info } from 'lucide-react'

export function DemoHint() {
  return (
    <p className="flex gap-3 rounded-xl bg-muted/60 p-4 text-sm text-muted-foreground">
      <Info className="mt-0.5 size-4 shrink-0" />
      <span>
        Just looking around? Sign in as <span className="font-medium text-foreground">demo@uptime.io</span> with
        password <span className="font-medium text-foreground">demo1234</span>.
      </span>
    </p>
  )
}
