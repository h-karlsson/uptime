import { delay, http } from 'msw'
import { authHandlers } from './auth'
import { healthHandlers } from './health'
import { monitorHandlers } from './monitors'
import { statsHandlers } from './stats'

export const handlers = [
  http.all('/api/*', () => delay()),
  ...healthHandlers,
  ...authHandlers,
  ...monitorHandlers,
  ...statsHandlers,
]
