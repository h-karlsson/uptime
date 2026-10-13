import type { ReactNode } from 'react'

type AuthStepProps = {
  title: string
  description: ReactNode
  children: ReactNode
}

export function AuthStep({ title, description, children }: AuthStepProps) {
  return (
    <div className="flex flex-col gap-7">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold tracking-tight">{title}</h1>
        <div className="text-muted-foreground">{description}</div>
      </div>
      {children}
    </div>
  )
}
