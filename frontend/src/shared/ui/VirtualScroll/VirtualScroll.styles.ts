import { styled } from '~/shared/styles'

export const triggerOffset = 500

export const StyledInnerDiv = styled('div', {
  position: 'relative',
  height: 'auto !important',
})

export const StyledTrigger = styled('span', {
  position: 'absolute',
  bottom: triggerOffset,
})
