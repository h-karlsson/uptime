import type { ReactElement, ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, renderHook } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TooltipProvider } from '@/components/ui/tooltip'

export function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
}

function wrapperFor(queryClient: QueryClient) {
  return ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>{children}</TooltipProvider>
    </QueryClientProvider>
  )
}

export function renderWithProviders(ui: ReactElement, queryClient = createTestQueryClient()) {
  return {
    user: userEvent.setup(),
    queryClient,
    ...render(ui, { wrapper: wrapperFor(queryClient) }),
  }
}

export function renderHookWithProviders<Result>(hook: () => Result, queryClient = createTestQueryClient()) {
  return {
    queryClient,
    ...renderHook(hook, { wrapper: wrapperFor(queryClient) }),
  }
}
