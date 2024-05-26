import { VariantProps } from '@stitches/react'
import React, { PropsWithChildren } from 'react'

import { StyledLayout, StyledRoot } from './Card.styles'

export interface CardProps extends PropsWithChildren, VariantProps<typeof StyledRoot> {
  padding?: number
}

export const Card: React.FC<CardProps> = ({ children, padding = 48, variant = 'elevate' }) => {
  return (
    <StyledRoot variant={variant} >
      <StyledLayout css={{ padding }}>
        {children}
      </StyledLayout>
    </StyledRoot>
  )
}
