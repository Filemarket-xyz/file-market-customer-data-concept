import { HTMLAttributes, Ref, RefObject } from 'react'
import { AriaButtonProps, mergeProps, useFocusRing, useHover, usePress } from 'react-aria'

import { useDOMRef } from '~/shared/lib'

export type AriaAnhorButtonProps = Omit<AriaButtonProps<'a'>, 'onFocus' | 'onBlur' | 'onKeyDown' | 'onKeyUp' | 'type' | 'target'>

export const useLink = <Props extends AriaAnhorButtonProps & HTMLAttributes<HTMLAnchorElement>, E extends HTMLElement>({
  isDisabled,
  onPress,
  ...props
}: Props, ref: RefObject<E | null> | Ref<E | null>) => {
  const linkRef = useDOMRef(ref)

  const { isHovered, hoverProps } = useHover({ isDisabled })
  const { isFocusVisible, focusProps } = useFocusRing()
  const { isPressed, pressProps } = usePress({ ref: linkRef, onPress })

  return {
    linkProps: {
      ...props,
      ...mergeProps(hoverProps, focusProps, pressProps),
      'data-hovered': isHovered,
      'data-focus-visible': isFocusVisible,
      'data-pressed': isPressed,
      'data-disabled': isDisabled,
    },
    linkRef,
  }
}
