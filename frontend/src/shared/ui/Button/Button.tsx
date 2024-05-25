import { ComponentProps, forwardRef } from 'react'
import { AriaButtonProps } from 'react-aria'

import { Loading } from '../Loading'
import { StyledButton } from './Button.styles'
import { useButton } from './useButton'

export interface ButtonProps extends AriaButtonProps, Omit<ComponentProps<typeof StyledButton>, 'onFocus' | 'onBlur' | 'onKeyDown' | 'onKeyUp' | 'isDisabled'> {
  loading?: boolean
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(({
  children,
  loading,
  isDisabled,
  ...props
}, ref) => {
  const { buttonRef, buttonProps } = useButton({
    ...props,
    loading,
    isDisabled: loading || isDisabled,
  }, ref)

  return (
    <StyledButton
      {...buttonProps}
      ref={buttonRef}
    >
      {children}
      <Loading
        color={props.variant === 'glowing' ? 'primary' : 'disabled'}
        loading={loading}
        size={20}
        fullWidth={false}
      />
    </StyledButton>
  )
})
