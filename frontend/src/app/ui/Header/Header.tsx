import React from 'react'

import { ConnectButton } from '~/features/Connect'
import { Container, Flex } from '~/shared/ui'

import { StyledHeader } from './Header.styles'
import { HeaderLogo } from './Logo'

export const Header: React.FC = () => {
  return (
    <StyledHeader>
      <Container fullHeight>
        <Flex
          fullHeight
          justifyContent='space-between'
          alignItems="center"
          flexWrap="nowrap"
          gap={30}
        >
          <HeaderLogo />
          <ConnectButton />
        </Flex>
      </Container>
    </StyledHeader>
  )
}
