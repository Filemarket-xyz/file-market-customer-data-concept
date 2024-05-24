import React from 'react'

import { Button, ButtonProps } from '../Button'
import { useButton } from '../useButton'
import { StyledGlow, StyledWrapper } from './ButtonGlowing.styles'

export const ButtonGlowing = React.forwardRef<HTMLButtonElement, Omit<ButtonProps, 'variant'>>(({
  children,
  ...props
}, ref) => {
  const { buttonRef, buttonProps } = useButton(props, ref)

  return (
    <StyledWrapper fullWidth={props.fullWidth}>
      <StyledGlow isDisabled={props.isDisabled} />
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
