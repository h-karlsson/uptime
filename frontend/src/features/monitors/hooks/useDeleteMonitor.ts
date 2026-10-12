import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deleteMonitor } from '../api'
import { monitorsQueryOptions } from '../queries'

export function useDeleteMonitor() {
  const queryClient = useQueryClient()
  const { queryKey } = monitorsQueryOptions

  return useMutation({
    mutationFn: deleteMonitor,
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey })
      const previous = queryClient.getQueryData(queryKey)
      queryClient.setQueryData(queryKey, (monitors) => monitors?.filter((monitor) => monitor.id !== id))
      return { previous }
    },
    onError: (_error, _id, context) => queryClient.setQueryData(queryKey, context?.previous),
  })
}
