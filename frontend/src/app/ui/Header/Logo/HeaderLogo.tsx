import React from 'react'
import { Link } from 'react-router-dom'

import { to } from '~/shared/lib'
import { Button, useButton } from '~/shared/ui'

import { Logo } from './HeaderLogo.styles'
import LogoIcon from './logo.svg'

export const HeaderLogo: React.FC = (() => {
  const {
    buttonRef,
    buttonProps,
  } = useButton<Record<any, never>, HTMLAnchorElement>({})

  return (
    <Link
      to={to.root()}
      {...buttonProps}
      ref={buttonRef}
    >
      <Button variant="text">
        <Logo
          src={LogoIcon}
          alt="FileMarket logo"
        />
      </Button>
    </Link>
  )
})
