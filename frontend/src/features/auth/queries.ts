import { queryOptions } from '@tanstack/react-query'
import { getMe } from './api'

export const meQueryOptions = queryOptions({
  queryKey: ['auth', 'me'],
  queryFn: ({ signal }) => getMe(signal),
  staleTime: Infinity,
})
