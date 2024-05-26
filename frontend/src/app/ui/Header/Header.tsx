import React from 'react'

import { ConnectButton } from '~/features/Auth'
import { Container, Flex } from '~/shared/ui'

import { Logo } from '../../../entities/Logo/ui'
import { StyledHeader } from './Header.styles'

export const Header: React.FC = () => {
  return (
    <StyledHeader>
      <Container fullHeight>
        <Flex
          fullHeight
          gap={16}
          alignItems="center"
          justifyContent="space-between"
          flexWrap="nowrap"
        >
          <Logo />
          <ConnectButton size="small" />
        </Flex>
      </Container>
    </StyledHeader>
  )
}
