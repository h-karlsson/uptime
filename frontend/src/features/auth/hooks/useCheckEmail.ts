import { useMutation } from '@tanstack/react-query'
import { checkEmail } from '../api'

export function useCheckEmail() {
  return useMutation({ mutationFn: checkEmail })
}
