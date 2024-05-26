import { VariantProps } from '@stitches/react'
import React, { PropsWithChildren } from 'react'

import { AppCSS } from '~/shared/styles'

import { StyledLayout, StyledRoot } from './Card.styles'

export interface CardProps extends PropsWithChildren, VariantProps<typeof StyledRoot>, Pick<AppCSS, 'padding'> {
}

export const Card: React.FC<CardProps> = ({ children, padding = 48, variant = 'outlined' }) => {
  return (
    <StyledRoot variant={variant} >
      <StyledLayout css={{ padding }}>
        {children}
      </StyledLayout>
    </StyledRoot>
  )
}
