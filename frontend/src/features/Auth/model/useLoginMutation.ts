import { useDisconnect } from 'thirdweb/react'
import { useSignMessage } from 'wagmi'

import { metamask } from '../lib'
import { useAuthBySignatureMutation } from './useAuthBySignatureMutation'
import { useAuthMessageMutation } from './useAuthMessageMutation'

export const useLoginMutation = () => {
  const { disconnect } = useDisconnect()
  const signMessage = useSignMessage({
    mutation: {
      onError: () => disconnect(metamask),
    },
  })

  const connectMessage = useAuthMessageMutation()
  const loginBySignature = useAuthBySignatureMutation()

  const mutateAsync = async () => {
    const response = await connectMessage.mutateAsync()
    const signature = await signMessage.signMessageAsync({ message: response.message })

    return loginBySignature.mutateAsync(signature)
  }

  return {
    isPending: loginBySignature.isPending || connectMessage.isPending || signMessage.isPending,
    mutateAsync,
    mutate: () => void mutateAsync(),
  }
}
