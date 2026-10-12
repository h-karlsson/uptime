import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createMonitor } from '../api'
import { monitorsQueryOptions } from '../queries'

export function useCreateMonitor() {
  const queryClient = useQueryClient()
  const { queryKey } = monitorsQueryOptions

  return useMutation({
    mutationFn: createMonitor,
    onSuccess: (created) => queryClient.setQueryData(queryKey, (monitors = []) => [...monitors, created]),
  })
}
