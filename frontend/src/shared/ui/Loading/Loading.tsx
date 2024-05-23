import React, { ComponentProps, type PropsWithChildren } from 'react'

import { Flex } from '../Flex'
import { StyledDiv } from './Loading.styles'

interface LoadingProps extends PropsWithChildren, ComponentProps<typeof StyledDiv> {
  loading?: boolean
  size?: number
}

export const Loading: React.FC<LoadingProps> = ({ loading, children, size = 24, ...props }) => {
  return (
    <>
      {loading ? (
        <Flex
          fullHeight
          fullWidth
          justifyContent='center'
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
