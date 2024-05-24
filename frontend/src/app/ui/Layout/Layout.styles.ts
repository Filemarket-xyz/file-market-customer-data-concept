import { styled } from '~/shared/styles'

export const StyledLayout = styled('main', {
  display: 'flex',
  flexDirection: 'column',
  width: '100%',
  height: '100%',
})

export const StyledBody = styled('div', {
  paddingTop: '$sizes$header',
  width: '100%',
  height: '100svh',
  minHeight: '100svh',
})
