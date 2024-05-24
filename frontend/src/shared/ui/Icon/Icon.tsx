import React from 'react'

import { AppCSS } from '~/shared/styles'

import { IconName, icons } from './icons'

export interface IconProps extends Pick<AppCSS, 'color'> {
  name: IconName
  className?: string
  size?: number | string | false
  css?: Omit<AppCSS, 'width' | 'height' | 'minHeight' | 'minWidth'>
}

export const Icon: React.FC<IconProps> = ({ name, className, size = 20, css, color }) => {
  const { Component } = icons[name]

  return (
    <Component
      className={className}
      css={{
        ...(size ? {
          width: size,
          height: size,
          minHeight: size,
          minWidth: size,
        } : undefined),
        '& > *': {
          color,
          fill: color,
        },
        ...css,
      }}
    />
  )
}
