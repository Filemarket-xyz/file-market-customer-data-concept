import { ConnectButton } from '~/features/Auth'
import { styled } from '~/shared/styles'
import { Container, Flex, Icon } from '~/shared/ui'

export const StyledContainer = styled(Container, {
  overflow: 'hidden',
  position: 'relative',
})

export const StyledWrapperFlex = styled(Flex, {
  height: 'calc(100vh - $sizes$header)',
  paddingTB: 20,
  justifyContent: 'center',
  gap: 100,
  '@xl': {
    gap: 20,
  },
  '@lg': {
    justifyContent: 'flex-start',
  },
})

export const StyledContentFlex = styled(Flex, {
  maxWidth: 640,
})

export const StyledConnectButton = styled(ConnectButton, {
  maxWidth: '70%',
  '@md': {
    maxWidth: 'none',
  },
})

export const StyledIcon = styled(Icon, {
  size: 800,
  '@lg': {
    zIndex: -1,
    position: 'absolute',
    size: 400,
    right: -150,
  },
  '@md': {
    display: 'none',
  },
})
