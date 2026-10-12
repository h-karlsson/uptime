import type { FormEvent } from 'react'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { SubmitButton } from '@/components/SubmitButton'
import { useLogin } from '../hooks/useLogin'
import { AuthCard } from './AuthCard'
import { EmailSummary } from './EmailSummary'

type PasswordStepProps = {
  email: string
  onBack: () => void
  onSuccess: () => void
}

export function PasswordStep({ email, onBack, onSuccess }: PasswordStepProps) {
  const login = useLogin()

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const password = String(new FormData(event.currentTarget).get('password'))
    login.mutate({ email, password }, { onSuccess })
  }

  return (
    <AuthCard title="Welcome back" description={<EmailSummary email={email} onChange={onBack} />}>
      <form onSubmit={handleSubmit}>
        {/* Lets password managers pair the saved password with this email. */}
        <input type="email" name="username" value={email} autoComplete="username" readOnly hidden />
        <FieldGroup>
          <Field data-invalid={login.isError}>
            <FieldLabel htmlFor="password">Password</FieldLabel>
            <Input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              aria-invalid={login.isError}
              required
              autoFocus
            />
            <FieldError>{login.error?.message}</FieldError>
          </Field>
          <SubmitButton pending={login.isPending}>Sign in</SubmitButton>
        </FieldGroup>
      </form>
    </AuthCard>
  )
}
