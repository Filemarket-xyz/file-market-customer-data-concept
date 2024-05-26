import React from 'react'

import { Flex, Txt } from '~/shared/ui'

import { StyledConnectButton, StyledContainer, StyledContentFlex, StyledIcon, StyledWrapperFlex } from './Main.page.styles'

export const MainPage: React.FC = () => {
  return (
    <StyledContainer>
      <StyledWrapperFlex alignItems="center">
        <StyledContentFlex flexDirection="column" gap={64}>
          <Flex flexDirection="column" gap={32}>
            <Txt fourfold1>
              Data Collection and Privacy protocol for ads and ML / AI
            </Txt>
            <Txt body1>
              Collaborative Data Collection protocol and platform that financiallycompensates users
              for sharing their data while providing companies with access to segmented audiences and hard-to-get data.
            </Txt>
          </Flex>
          <StyledConnectButton fullWidth glowing />
        </StyledContentFlex>
        <StyledIcon name="storage" size={false} />
      </StyledWrapperFlex>
    </StyledContainer>
  )
}
