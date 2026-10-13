import { screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { authHandlers } from '@/test/handlers/auth'
import { renderWithProviders } from '@/test/render'
import { server } from '@/test/server'
import { LoginFlow } from './LoginFlow'

function setup() {
  const onAuthenticated = vi.fn()
  const rendered = renderWithProviders(<LoginFlow onAuthenticated={onAuthenticated} />)
  const continueWithEmail = async (email: string) => {
    await rendered.user.type(screen.getByLabelText('Email'), email)
    await rendered.user.click(screen.getByRole('button', { name: 'Continue' }))
  }
  return { ...rendered, onAuthenticated, continueWithEmail }
}

describe('LoginFlow', () => {
  beforeEach(() => server.use(...authHandlers))

  it('asks a registered user for their password and signs them in', async () => {
    const { user, onAuthenticated, continueWithEmail } = setup()

    await continueWithEmail('demo@uptime.io')
    expect(await screen.findByText('Welcome back')).toBeInTheDocument()

    await user.type(screen.getByLabelText('Password'), 'demo1234')
    await user.click(screen.getByRole('button', { name: 'Sign in' }))

    await vi.waitFor(() => expect(onAuthenticated).toHaveBeenCalled())
  })

  it('shows the server error for a wrong password', async () => {
    const { user, onAuthenticated, continueWithEmail } = setup()

    await continueWithEmail('demo@uptime.io')
    await user.type(await screen.findByLabelText('Password'), 'wrongpass')
    await user.click(screen.getByRole('button', { name: 'Sign in' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('invalid email or password')
    expect(onAuthenticated).not.toHaveBeenCalled()
  })

  it('offers registration for a new email and checks the passwords match', async () => {
    const { user, onAuthenticated, continueWithEmail } = setup()

    await continueWithEmail('new@uptime.io')
    expect(await screen.findByText('Create your account')).toBeInTheDocument()

    await user.type(screen.getByLabelText('Password'), 'longenough')
    await user.type(screen.getByLabelText('Confirm password'), 'different')
    await user.click(screen.getByRole('button', { name: 'Create account' }))
    expect(await screen.findByText("Passwords don't match.")).toBeInTheDocument()

    await user.clear(screen.getByLabelText('Confirm password'))
    await user.type(screen.getByLabelText('Confirm password'), 'longenough')
    await user.click(screen.getByRole('button', { name: 'Create account' }))

    await vi.waitFor(() => expect(onAuthenticated).toHaveBeenCalled())
  })

  it('lets the user go back and change the email', async () => {
    const { user, continueWithEmail } = setup()

    await continueWithEmail('demo@uptime.io')
    await user.click(await screen.findByRole('button', { name: 'Change' }))

    expect(screen.getByLabelText('Email')).toHaveValue('demo@uptime.io')
  })
})
