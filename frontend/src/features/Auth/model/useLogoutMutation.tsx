import { useMutation, useQueryClient } from '@tanstack/react-query'

import { userKeys } from '~/entities/User'
import { api } from '~/shared/api'

export const useLogoutMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: userKeys.logout(),
    mutationFn: () => api.auth.logoutCreate(),
    onSettled: () => {
      queryClient.setQueryData(userKeys.root, null)
    },
  })
}
