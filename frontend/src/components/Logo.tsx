import { Activity } from 'lucide-react'
import { cn } from '@/lib/utils'

export function Logo({ className }: { className?: string }) {
  return (
    <div className={cn('flex items-center gap-2 font-semibold', className)}>
      <div className="flex size-7 items-center justify-center rounded-md bg-primary text-primary-foreground">
        <Activity className="size-4" />
      </div>
      Uptime
    </div>
  )
}
