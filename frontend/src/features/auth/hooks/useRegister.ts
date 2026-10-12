import { useMutation, useQueryClient } from '@tanstack/react-query'
import { register } from '../api'
import { meQueryOptions } from '../queries'

export function useRegister() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: register,
    onSuccess: (user) => queryClient.setQueryData(meQueryOptions.queryKey, user),
  })
}
