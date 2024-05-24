import React from 'react'

import { ConnectButton } from '~/features/Connect'
import { Container, Flex } from '~/shared/ui'

import { Logo } from '../../../entities/Logo/ui'
import { StyledHeader } from './Header.styles'

export const Header: React.FC = () => {
  return (
    <StyledHeader>
      <Container fullHeight>
        <Flex
          fullHeight
          gap={30}
          alignItems="center"
          justifyContent="space-between"
          flexWrap="nowrap"
        >
          <Logo />
          <ConnectButton />
        </Flex>
      </Container>
    </StyledHeader>
  )
}
