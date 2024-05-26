import React from 'react'

import { Button, ButtonProps } from '../Button'
import { StyledGlow, StyledWrapper } from './ButtonGlowing.styles'

export interface ButtonGlowingProps extends Omit<ButtonProps, 'variant'> {}

export const ButtonGlowing = React.forwardRef<HTMLButtonElement, ButtonGlowingProps>(({
  children,
  className,
  isDisabled,
  loading,
  ...props
}, ref) => {
  return (
    <StyledWrapper fullWidth={props.fullWidth} className={className}>
      <StyledGlow isDisabled={isDisabled || loading} />
      <Button
        ref={ref}
        {...props}
        loading={loading}
        isDisabled={isDisabled || loading}
        variant="glowing"
      >
        {children}
      </Button>
    </StyledWrapper>
  )
})
