import { useEffect } from 'react'
import { useQueryErrorResetBoundary } from '@tanstack/react-query'
import { useRouter, type ErrorComponentProps } from '@tanstack/react-router'
import { CircleAlert, RotateCw } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function RouteError({ error }: ErrorComponentProps) {
  const router = useRouter()
  const queryErrorResetBoundary = useQueryErrorResetBoundary()

  useEffect(() => {
    queryErrorResetBoundary.reset()
  }, [queryErrorResetBoundary])

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-3 p-6 text-center">
      <CircleAlert className="size-8 text-muted-foreground" />
      <p className="font-medium">Something went wrong</p>
      <p className="max-w-sm text-sm text-muted-foreground">
        {error instanceof Error ? error.message : "The page couldn't be loaded."}
      </p>
      <Button variant="outline" onClick={() => router.invalidate()}>
        <RotateCw />
        Try again
      </Button>
    </div>
  )
}
