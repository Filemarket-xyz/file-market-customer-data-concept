import { styled } from '~/shared/styles'
import { Flex } from '~/shared/ui'

export const StyledFlex = styled(Flex, {
  '@sm': {
    flexFlow: 'column-reverse wrap',
    alignItems: 'center',
    gap: '$3',
  },
})

export const StyledHr = styled('hr', {
  border: 'none',
  width: '1px',
  height: '18px',
  background: '#232528',
  borderRadius: '2px',
  '@sm': {
    display: 'none',
  },
})
