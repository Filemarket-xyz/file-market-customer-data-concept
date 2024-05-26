import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useActiveAccount } from 'thirdweb/react'

import { userKeys } from '~/entities/User'
import { AuthResponse, UserRole } from '~/swagger/Api'

export const useAuthBySignatureMutation = () => {
  const queryClient = useQueryClient()
  const account = useActiveAccount()
  const address = account?.address ?? ''

  return useMutation({
    mutationKey: userKeys.loginBySignature(address),
    // mutationFn: (signature: string) => api.auth.bySignatureCreate({ address, signature }).then(({ data }) => data),
    mutationFn: (signature: string) => {
      return new Promise<AuthResponse>((resolve) => resolve({ user: { role: UserRole.Client } }))
    },
    onSuccess: (data) => {
      queryClient.setQueryData<AuthResponse>(userKeys.root, data)
    },
  })

}
