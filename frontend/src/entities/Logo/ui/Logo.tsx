import React from 'react'

import { to } from '~/shared/lib'
import { Icon, IconProps, Link } from '~/shared/ui'

export interface LogoProps extends Pick<IconProps, 'color'> {}

export const Logo: React.FC<LogoProps> = (({ color }) => {
  return (
    <Link to={to.root()}>
      <Icon size={false} name="logo" color={color} />
    </Link>
  )
})
