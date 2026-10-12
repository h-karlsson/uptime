import { useState, type FormEvent } from 'react'
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { useRegister } from '../hooks/useRegister'
import { AuthCard } from './AuthCard'
import { EmailSummary } from './EmailSummary'
import { SubmitButton } from './SubmitButton'

// Matches the backend rule; the backend still validates, this just gives earlier feedback.
const MIN_PASSWORD_LENGTH = 8

type RegisterStepProps = {
  email: string
  onBack: () => void
  onSuccess: () => void
}

export function RegisterStep({ email, onBack, onSuccess }: RegisterStepProps) {
  const register = useRegister()
  const [passwordsMatch, setPasswordsMatch] = useState(true)

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    const password = String(form.get('password'))
    const matches = password === form.get('confirmPassword')
    setPasswordsMatch(matches)
    if (matches) {
      register.mutate({ email, password }, { onSuccess })
    }
  }

  return (
    <AuthCard
      title="Create your account"
      description={
        <>
          No account yet for <EmailSummary email={email} onChange={onBack} />
        </>
      }
    >
      <form onSubmit={handleSubmit}>
        <input type="email" name="username" value={email} autoComplete="username" readOnly hidden />
        <FieldGroup>
          <Field data-invalid={register.isError}>
            <FieldLabel htmlFor="password">Password</FieldLabel>
            <Input
              id="password"
              name="password"
              type="password"
              autoComplete="new-password"
              minLength={MIN_PASSWORD_LENGTH}
              aria-invalid={register.isError}
              required
              autoFocus
            />
            {register.isError ? (
              <FieldError>{register.error.message}</FieldError>
            ) : (
              <FieldDescription>At least {MIN_PASSWORD_LENGTH} characters.</FieldDescription>
            )}
          </Field>
          <Field data-invalid={!passwordsMatch}>
            <FieldLabel htmlFor="confirmPassword">Confirm password</FieldLabel>
            <Input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              autoComplete="new-password"
              aria-invalid={!passwordsMatch}
              required
            />
            {!passwordsMatch && <FieldError>Passwords don't match.</FieldError>}
          </Field>
          <SubmitButton pending={register.isPending}>Create account</SubmitButton>
        </FieldGroup>
      </form>
    </AuthCard>
  )
}
