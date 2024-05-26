import React from 'react'
import { useActiveAccount, useConnect, useDisconnect } from 'thirdweb/react'

import { client } from '~/shared/lib'
import { Button, ButtonGlowing, ButtonProps } from '~/shared/ui'

import { metamask } from '../../lib/wallets'
import { useLoginMutation, useLogoutMutation } from '../../model'

export interface ConnectButtonProps extends Omit<ButtonProps, 'variant' | 'onPress' | 'loading'> {
  glowing?: boolean
}

export const ConnectButton: React.FC<ConnectButtonProps> = ({ glowing, ...props }) => {
  const ButtonComponent = glowing ? ButtonGlowing : Button

  const { connect, isConnecting } = useConnect()
  const { disconnect } = useDisconnect()
  const account = useActiveAccount()
  const login = useLoginMutation()
  const logout = useLogoutMutation()

  const connectWallet = async () => {
    await connect(async() => {
      await metamask.connect({ client })

      return metamask
    })

    login.mutate()
  }

  const disconnectWallet = () => {
    disconnect(metamask)
    logout.mutate()
  }

  const onPress = () => {
    if (account) {
      disconnectWallet()
    } else {
      connectWallet()
    }
  }

  return (
    <ButtonComponent
      {...props}
      variant="secondary"
      loading={isConnecting || login.isPending}
      onPress={onPress}
    >
      {account && !login.isPending ? 'Disconnect' : 'Connect'}
    </ButtonComponent>
  )
}
