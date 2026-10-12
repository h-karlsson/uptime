import { useState } from 'react'
import { EmailStep } from './EmailStep'
import { PasswordStep } from './PasswordStep'
import { RegisterStep } from './RegisterStep'

type Step = 'email' | 'password' | 'register'

/**
 * Email-first sign in: ask for the email, then show the password field for existing
 * accounts or the registration form for new ones.
 */
export function LoginFlow({ onAuthenticated }: { onAuthenticated: () => void }) {
  const [step, setStep] = useState<Step>('email')
  const [email, setEmail] = useState('')

  function handleEmail(email: string, registered: boolean) {
    setEmail(email)
    setStep(registered ? 'password' : 'register')
  }

  const backToEmail = () => setStep('email')

  switch (step) {
    case 'email':
      return <EmailStep defaultEmail={email} onContinue={handleEmail} />
    case 'password':
      return <PasswordStep email={email} onBack={backToEmail} onSuccess={onAuthenticated} />
    case 'register':
      return <RegisterStep email={email} onBack={backToEmail} onSuccess={onAuthenticated} />
  }
}
