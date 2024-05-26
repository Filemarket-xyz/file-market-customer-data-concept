import React, { PropsWithChildren } from 'react'

import { Flex } from '../../Flex'
import { Icon, IconName } from '../../Icon'
import { Txt } from '../../Txt'
import { Card, CardProps } from '../Card'

export interface ActionCardProps extends PropsWithChildren, Pick<CardProps, 'variant'> {
  iconName?: IconName
  title: string
  description?: string
}

export const ActionCard: React.FC<ActionCardProps> = ({
  children,
  iconName,
  title,
  description,
  variant = 'outline',
}) => {
  return (
    <Card variant={variant} padding={32}>
      <Flex gap={40} alignItems="center">
        {iconName && <Icon color="$blue500" name={iconName} size={110} />}
        <Flex
          fullWidth
          fullHeight
          gap={12}
          flexDirection="column"
          justifyContent="space-between"
        >
          <Flex fullWidth gap={12} flexDirection="column">
            <Txt primary1 color="gray800" css={{ fontSize: 24 }}>
              {title}
            </Txt>
            {description && (
              <Txt body4 color="gray700">
                {description}
              </Txt>
            )}
          </Flex>
          {children}
        </Flex>
      </Flex>
    </Card>
  )
}
