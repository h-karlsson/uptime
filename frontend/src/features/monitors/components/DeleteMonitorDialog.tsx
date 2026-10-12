import type { Monitor } from '@/client'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { useDeleteMonitor } from '../hooks/useDeleteMonitor'

type DeleteMonitorDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  monitor?: Monitor
}

export function DeleteMonitorDialog({ open, onOpenChange, monitor }: DeleteMonitorDialogProps) {
  const deleteMonitor = useDeleteMonitor()

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete {monitor?.name}?</AlertDialogTitle>
          <AlertDialogDescription>
            The monitor and its check history are deleted permanently.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction variant="destructive" onClick={() => monitor && deleteMonitor.mutate(monitor.id)}>
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
