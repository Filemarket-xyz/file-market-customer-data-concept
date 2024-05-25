import React, { ComponentProps, type PropsWithChildren } from 'react'

import { Flex, FlexProps } from '../Flex'
import { StyledDiv } from './Loading.styles'

interface LoadingProps extends PropsWithChildren, ComponentProps<typeof StyledDiv>, Pick<FlexProps, 'fullWidth' | 'fullHeight'> {
  loading?: boolean
  size?: number
}

export const Loading: React.FC<LoadingProps> = ({
  loading,
  children,
  size = 24,
  fullHeight = true,
  fullWidth = true,
  ...props
}) => {
  return (
    <>
      {loading ? (
        <Flex
          fullHeight={fullHeight}
          fullWidth={fullWidth}
          justifyContent="center"
          alignItems="center"
          css={{ minHeight: size }}
        >
          <StyledDiv
            {...props}
            css={{
              size,
              ...props.css,
            }}
          />
        </Flex>
      ) : (
        children
      )}
    </>
  )
}
