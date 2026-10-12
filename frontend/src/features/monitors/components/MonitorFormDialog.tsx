import type { Monitor } from '@/client'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { MonitorForm } from './MonitorForm'

type MonitorFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  monitor?: Monitor
}

export function MonitorFormDialog({ open, onOpenChange, monitor }: MonitorFormDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[calc(100svh-2rem)] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{monitor ? 'Edit monitor' : 'Add monitor'}</DialogTitle>
          <DialogDescription>We'll request the URL on every interval and alert you when it goes down.</DialogDescription>
        </DialogHeader>
        <MonitorForm monitor={monitor} onDone={() => onOpenChange(false)} />
      </DialogContent>
    </Dialog>
  )
}
