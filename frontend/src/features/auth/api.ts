import * as sdk from '@/client'
import type { Credentials, User } from '@/client'
import { ApiError } from '@/lib/api'

export async function getMe(signal?: AbortSignal): Promise<User | null> {
  try {
    const { data } = await sdk.getMe({ signal })
    return data
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null
    }
    throw error
  }
}

export async function checkEmail(email: string) {
  const { data } = await sdk.checkEmail({ body: { email } })
  return data
}

export async function login(credentials: Credentials) {
  const { data } = await sdk.login({ body: credentials })
  return data
}

export async function register(credentials: Credentials) {
  const { data } = await sdk.register({ body: credentials })
  return data
}

export async function logout() {
  await sdk.logout()
}
