import React, { PropsWithChildren } from 'react'

import { AppCSS } from '~/shared/styles'

import { StyledDiv } from './Flex.styles'

interface FlexProps extends PropsWithChildren, Pick<AppCSS, 'gap' | 'alignItems' | 'justifyContent' | 'flexDirection' | 'flexWrap' | 'flexGrow'> {
  css?: AppCSS
  fullWidth?: boolean
  fullHeight?: boolean
  className?: string
}

export const Flex: React.FC<FlexProps> = ({ children, fullHeight, fullWidth, css, className, ...props }) => {
  return (
    <StyledDiv
      fullHeight={fullHeight}
      fullWidth={fullWidth}
      className={className}
      css={{
        ...props,
        ...css,
      }}
    >
      {children}
    </StyledDiv>
  )
}
