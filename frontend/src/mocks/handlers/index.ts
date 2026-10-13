import { delay, http } from 'msw'
import { authHandlers } from './auth'
import { monitorHandlers } from './monitors'
import { statsHandlers } from './stats'

export const handlers = [
  // Returns nothing, so every request continues to its real handler after a realistic network delay.
  http.all('/api/*', () => delay()),
  ...authHandlers,
  ...monitorHandlers,
  ...statsHandlers,
]
