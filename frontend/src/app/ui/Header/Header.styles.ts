import { styled } from '~/shared/styles'

export const StyledHeader = styled('header', {
  width: '100%',
  position: 'fixed',
  zIndex: 10,
  top: 0,
  left: 0,
  right: 0,
  boxShadow: '$header',
  color: '$gray600',
  background: 'rgba(249, 249, 249, 0.75)',
  backdropFilter: 'blur(14px)',
  transition: 'all 0.5s',
  height: 80,
})
