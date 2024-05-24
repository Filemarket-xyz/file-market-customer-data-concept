import React from 'react'
import { useActiveAccount, useConnect, useDisconnect } from 'thirdweb/react'
import { createWallet } from 'thirdweb/wallets'
import { useSignMessage } from 'wagmi'

import { client } from '~/shared/lib'
import { Button } from '~/shared/ui'

const metamask = createWallet('io.metamask')

export const ConnectButton: React.FC = () => {
  const { connect, isConnecting } = useConnect()
  const { disconnect } = useDisconnect()
  const { signMessageAsync } = useSignMessage()
  const account = useActiveAccount()

  const onPress = () => {
    if (account) {
      disconnect(metamask)
    } else {
      connect(async() => {
        await metamask.connect({ client })

        await signMessageAsync({ message: 'Hello, I am using WalletConnect' })

        return metamask
      })
    }
  }

  return (
    <Button
      variant="secondary"
      size="small"
      loading={isConnecting}
      onPress={onPress}
    >
      {account ? 'Disconnect' : 'Connect'}
    </Button>
  )
}
