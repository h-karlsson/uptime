import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { Logo } from '@/components/Logo'
import { ThemeToggle } from '@/components/ThemeToggle'
import { LoginFlow } from '@/features/auth/components/LoginFlow'
import { meQueryOptions } from '@/features/auth/queries'
import { safeRedirect } from '@/features/auth/redirect'

type LoginSearch = {
  redirect?: string
}

export const Route = createFileRoute('/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => ({
    redirect: typeof search.redirect === 'string' ? search.redirect : undefined,
  }),
  // Already logged in? Skip the login page.
  beforeLoad: async ({ context, search }) => {
    const user = await context.queryClient.ensureQueryData(meQueryOptions)
    if (user) {
      throw redirect({ href: safeRedirect(search.redirect) })
    }
  },
  component: LoginPage,
})

function LoginPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()

  return (
    <div className="relative flex min-h-svh flex-col items-center justify-center bg-muted p-6 md:p-10">
      <div className="absolute top-4 right-4">
        <ThemeToggle />
      </div>
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Logo className="self-center" />
        <LoginFlow onAuthenticated={() => void navigate({ href: safeRedirect(search.redirect) })} />
        <p className="text-center text-sm text-balance text-muted-foreground">
          Just looking around? Sign in as <span className="font-medium text-foreground">demo@uptime.io</span>{' '}
          with password <span className="font-medium text-foreground">demo1234</span>.
        </p>
      </div>
    </div>
  )
}
