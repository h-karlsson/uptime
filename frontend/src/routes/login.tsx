import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { Logo } from '@/components/Logo'
import { ThemeToggle } from '@/components/ThemeToggle'
import { DemoHint } from '@/features/auth/components/DemoHint'
import { LoginFlow } from '@/features/auth/components/LoginFlow'
import { LoginShowcase } from '@/features/auth/components/LoginShowcase'
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
    <div className="grid min-h-svh bg-background lg:grid-cols-2">
      <div className="flex min-w-0 flex-col gap-8 p-6 md:px-8">
        <div className="flex items-center justify-between">
          <Logo />
          <ThemeToggle />
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="flex w-full max-w-sm flex-col gap-7">
            <LoginFlow onAuthenticated={() => void navigate({ href: safeRedirect(search.redirect) })} />
            <DemoHint />
          </div>
        </div>
      </div>
      <LoginShowcase />
    </div>
  )
}
