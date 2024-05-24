import { HTMLAttributes, MouseEventHandler, Ref, RefObject, useCallback } from 'react'
import { AriaButtonProps, mergeProps, useButton as useButtonAria, useFocusRing, useHover } from 'react-aria'

import { useDOMRef } from '~/shared/lib'

export const useButton = <Props extends AriaButtonProps<'button'> & HTMLAttributes<HTMLSpanElement>, E extends HTMLElement>({
  isDisabled,
  onPress,
  onPressStart,
  onPressEnd,
  onPressChange,
  onPressUp,
  onClick: deprecatedOnClick,
  ...props
}: Props, ref?: RefObject<E | null> | Ref<E | null>) => {
  const buttonRef = useDOMRef(ref)

  const { isPressed, buttonProps } = useButtonAria({
    isDisabled,
    onPress,
    onPressStart,
    onPressEnd,
    onPressChange,
    onPressUp,
  }, buttonRef)

  const onClick = useCallback<MouseEventHandler<HTMLButtonElement>>(
    (event) => {
      buttonProps.onClick?.(event)
      deprecatedOnClick?.(event)
    },
    [buttonProps.onClick],
  )

  const { isFocusVisible, focusProps } = useFocusRing()
  const { hoverProps, isHovered } = useHover({ isDisabled })

  return {
    buttonRef,
    buttonProps: {
      ...mergeProps(buttonProps, focusProps, hoverProps),
      ...props,
      onClick,
      'data-pressed': isPressed,
      'data-hovered': isHovered,
      'data-focus-ring': isFocusVisible,
      'data-disabled': isDisabled,
    },
  }
}
