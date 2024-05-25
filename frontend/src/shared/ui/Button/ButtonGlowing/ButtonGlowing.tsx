import React from 'react'

import { Button, ButtonProps } from '../Button'
import { useButton } from '../useButton'
import { StyledGlow, StyledWrapper } from './ButtonGlowing.styles'

export const ButtonGlowing = React.forwardRef<HTMLButtonElement, Omit<ButtonProps, 'variant'>>(({
  children,
  className,
  isDisabled,
  loading,
  ...props
}, ref) => {
  const { buttonRef, buttonProps } = useButton({
    ...props,
    loading,
    isDisabled: isDisabled || loading,
  }, ref)

  return (
    <StyledWrapper fullWidth={props.fullWidth} className={className}>
      <StyledGlow isDisabled={isDisabled || loading} />
      <Button
        {...buttonProps}
        ref={buttonRef}
        variant="glowing"
      >
        {children}
      </Button>
    </StyledWrapper>
  )
})
