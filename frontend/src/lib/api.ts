import { client } from '@/client/client.gen'

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export function setupApiClient() {
  client.interceptors.error.use((error, response) => {
    const status = response?.status ?? 0
    return new ApiError(status, errorMessage(error, status))
  })
}

function errorMessage(error: unknown, status: number): string {
  if (typeof error === 'object' && error !== null && 'error' in error && typeof error.error === 'string') {
    return error.error
  }
  return status ? `Request failed (${status})` : 'Could not reach the server'
}
