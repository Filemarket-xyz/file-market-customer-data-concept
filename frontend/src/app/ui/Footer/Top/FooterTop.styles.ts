import { styled } from '~/shared/styles'
import { Flex, Link, textVariant } from '~/shared/ui'

export const StyledWrapperFlex = styled(Flex, {
  gap: 16,
  justifyContent: 'space-between',
  '@sm': {
    justifyContent: 'flex-start',
    gap: '32px',
    columnGap: '74px',
  },
})

export const StyledFirstColumnFlex = styled(Flex, {
  maxWidth: '360px',
  '@md': {
    maxWidth: '100%',
    width: '100%',
  },
})

export const StyledMiddleColumnFlex = styled(Flex, {
  '@lg': {
    maxWidth: '360px',
  },
  '@sm': {
    maxWidth: '140px',
  },
})

export const StyledLastColumnFlex = styled(Flex, {
  maxWidth: '256px',
  '@md': {
    maxWidth: '100%',
    width: '100%',
  },
})

export const StyledH4 = styled('h4', {
  ...textVariant('secondary2').true,
  fontFamily: '$body',
  fontWeight: '700',
  color: '#7B7C7E',
})

export const StyledLink = styled(Link, {
  height: '44px',
  display: 'flex',
  flex: '1 0',
  alignItems: 'center',
  gap: '4px',
  background: '$gray800',
  borderRadius: '8px',
  paddingLR: 16,
  '@lg': {
    justifyContent: 'center',
    flex: '1 0 33%',
  },
  '&[data-hovered=true]': {
    opacity: 1,
    background: '#393B3E',
  },
})
