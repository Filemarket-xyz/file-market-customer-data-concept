import React, { PropsWithChildren, ReactNode } from 'react'

import { Flex } from '../../Flex'
import { Txt } from '../../Txt'

export interface CardTitleProps extends PropsWithChildren {
  prefix?: ReactNode
}

export const CardTitle: React.FC<CardTitleProps> = ({ prefix, children }) => {
  return (
    <Flex gap={12} alignItems="center">
      {prefix}
      <Txt fourfold1>
        {children}
      </Txt>
    </Flex>
  )
}
