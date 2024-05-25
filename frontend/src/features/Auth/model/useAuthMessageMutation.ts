import { useMutation } from '@tanstack/react-query'
import { useActiveAccount } from 'thirdweb/react'

import { userKeys } from '~/entities/User'
import { AuthMessageResponse } from '~/swagger/Api'

export const useAuthMessageMutation = () => {
  const account = useActiveAccount()
  const address = account?.address ?? ''

  return useMutation({
    mutationKey: userKeys.connectMessage(address),
    // mutationFn: () => api.auth.messageCreate({ address }).then(({ data }) => data),
    mutationFn: () => new Promise<AuthMessageResponse>((resolve) => resolve({ message: address })),
  })
}
