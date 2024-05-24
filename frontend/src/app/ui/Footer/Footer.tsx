import React from 'react'

import { FooterBottom } from './Bottom'
import { StyledContainer, StyledFooter, StyledHr } from './Footer.styles'
import { FooterTop } from './Top'

export const Footer: React.FC = () => {
  return (
    <StyledFooter>
      <StyledContainer>
        <FooterTop />
        <StyledHr />
        <FooterBottom />
      </StyledContainer>
    </StyledFooter>
  )
}
