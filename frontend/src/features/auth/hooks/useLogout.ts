import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { logout } from '../api'
import { meQueryOptions } from '../queries'

export function useLogout() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  return useMutation({
    mutationFn: logout,
    onSuccess: async () => {
      queryClient.clear()
      queryClient.setQueryData(meQueryOptions.queryKey, null)
      await navigate({ to: '/login' })
    },
  })
}
