import { VariantProps } from '@stitches/react'
import { forwardRef } from 'react'
import { NavLinkProps, To } from 'react-router-dom'

import { StyledAnhor, StyledNavLink } from './Link.styles'
import { AriaAnhorButtonProps, useLink } from './useLink'

export interface LinkProps extends Omit<NavLinkProps, 'children' | 'className' | 'style' | 'color' | 'to'>, AriaAnhorButtonProps, VariantProps<typeof StyledAnhor> {
  to?: To
}

export const Link = forwardRef<HTMLAnchorElement, LinkProps>(({
  // to means internal link
  to,
  // href means external link
  href,
  color,
  ...props
}, ref) => {
  const { linkRef, linkProps } = useLink(props, ref)

  if (!to && !href) {
    throw new Error('<Link /> must have either `to` or `href`')
  }

  return (
    <>
      {href && (
        <StyledAnhor
          ref={linkRef}
          href={href}
          color={color}
          {...linkProps}
        />
      )}

      {to && (
        <StyledNavLink
          ref={linkRef}
          to={to}
          color={color}
          {...linkProps}
        />
      )}
    </>
  )
})
