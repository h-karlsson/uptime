import { useMutation, useQueryClient } from '@tanstack/react-query'
import { login } from '../api'
import { meQueryOptions } from '../queries'

export function useLogin() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: login,
    onSuccess: (user) => queryClient.setQueryData(meQueryOptions.queryKey, user),
  })
}
