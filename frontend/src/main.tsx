import { StrictMode } from 'react'
import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { createRoot } from 'react-dom/client'
import { meQueryOptions } from '@/features/auth/queries'
import { ApiError, setupApiClient } from '@/lib/api'
import { syncThemeWithDocument } from '@/state/useThemeStore'
import { routeTree } from './routeTree.gen'
import './index.css'

setupApiClient()
syncThemeWithDocument()

const queryClient = new QueryClient({
  queryCache: new QueryCache({ onError: handleUnauthorized }),
  mutationCache: new MutationCache({ onError: handleUnauthorized }),
})

const router = createRouter({
  routeTree,
  context: { queryClient },
  defaultPreload: 'intent',
  scrollRestoration: true,
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

// A 401 from any request while we think someone is logged in means the session expired:
// forget the user and re-run the route guards, which send them to /login.
function handleUnauthorized(error: Error) {
  const loggedIn = queryClient.getQueryData(meQueryOptions.queryKey)
  if (error instanceof ApiError && error.status === 401 && loggedIn) {
    queryClient.setQueryData(meQueryOptions.queryKey, null)
    void router.invalidate()
  }
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
