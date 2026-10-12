import { queryOptions } from '@tanstack/react-query'
import { listMonitors } from './api'

export const monitorsQueryOptions = queryOptions({
  queryKey: ['monitors'],
  queryFn: ({ signal }) => listMonitors(signal),
})
