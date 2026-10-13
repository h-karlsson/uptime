import type { FormEvent } from 'react'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { SubmitButton } from '@/components/SubmitButton'
import { useCheckEmail } from '../hooks/useCheckEmail'
import { AuthStep } from './AuthStep'

type EmailStepProps = {
  defaultEmail: string
  onContinue: (email: string, registered: boolean) => void
}

export function EmailStep({ defaultEmail, onContinue }: EmailStepProps) {
  const checkEmail = useCheckEmail()

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const email = String(new FormData(event.currentTarget).get('email')).trim()
    checkEmail.mutate(email, {
      onSuccess: ({ registered }) => onContinue(email, registered),
    })
  }

  return (
    <AuthStep title="Sign in to Uptime" description="Enter your email to sign in or create an account.">
      <form onSubmit={handleSubmit}>
        <FieldGroup>
          <Field data-invalid={checkEmail.isError}>
            <FieldLabel htmlFor="email">Email</FieldLabel>
            <Input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              defaultValue={defaultEmail}
              aria-invalid={checkEmail.isError}
              required
              autoFocus
            />
            <FieldError>{checkEmail.error?.message}</FieldError>
          </Field>
          <SubmitButton pending={checkEmail.isPending}>Continue</SubmitButton>
        </FieldGroup>
      </form>
    </AuthStep>
  )
}
