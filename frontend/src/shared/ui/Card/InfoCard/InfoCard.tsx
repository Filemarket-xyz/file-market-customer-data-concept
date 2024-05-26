import React, { ReactNode } from 'react'

import { Flex } from '../../Flex'
import { Txt } from '../../Txt'
import { Card } from '../Card'
import { CardTitle } from '../Title'

export interface InfoCardProps {
  titlePrefix?: ReactNode
  title: string
  description: string
}

export const InfoCard: React.FC<InfoCardProps> = ({ titlePrefix, title, description }) => {
  return (
    <Card variant="elevated">
      <Flex gap={32} flexDirection="column" alignItems="center">
        <CardTitle prefix={titlePrefix}>
          {title}
        </CardTitle>
        <Txt body1 color="gray700" css={{ textAlign: 'center' }}>
          {description}
        </Txt>
      </Flex>
    </Card>
  )
}
