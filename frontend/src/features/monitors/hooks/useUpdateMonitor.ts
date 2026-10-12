import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { MonitorInput } from '@/client'
import { updateMonitor } from '../api'
import { monitorsQueryOptions } from '../queries'

type UpdateVariables = {
  id: number
  input: MonitorInput
}

// Optimistic, so pausing and resuming feel instant; the change is rolled back if the request fails.
export function useUpdateMonitor() {
  const queryClient = useQueryClient()
  const { queryKey } = monitorsQueryOptions

  return useMutation({
    mutationFn: ({ id, input }: UpdateVariables) => updateMonitor(id, input),
    onMutate: async ({ id, input }) => {
      await queryClient.cancelQueries({ queryKey })
      const previous = queryClient.getQueryData(queryKey)
      queryClient.setQueryData(queryKey, (monitors) =>
        monitors?.map((monitor) => (monitor.id === id ? { ...monitor, ...input } : monitor)),
      )
      return { previous }
    },
    onError: (_error, _variables, context) => queryClient.setQueryData(queryKey, context?.previous),
    onSuccess: (updated) =>
      queryClient.setQueryData(queryKey, (monitors) =>
        monitors?.map((monitor) => (monitor.id === updated.id ? updated : monitor)),
      ),
  })
}
