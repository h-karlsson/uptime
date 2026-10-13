import type { QueryClient } from '@tanstack/react-query'
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { MockModeBadge } from '@/components/MockModeBadge'
import { TooltipProvider } from '@/components/ui/tooltip'

export type RouterContext = {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
})

function RootLayout() {
  return (
    <TooltipProvider>
      <Outlet />
      {import.meta.env.MODE === 'mock' && <MockModeBadge />}
    </TooltipProvider>
  )
}
