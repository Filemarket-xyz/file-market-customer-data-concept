import { styled } from '~/shared/styles'

export const headerHeight = 80

export const StyledHeader = styled('header', {
  width: '100%',
  position: 'fixed',
  zIndex: 10,
  top: 0,
  left: 0,
  right: 0,
  boxShadow: '0px 4px 15px rgba(19, 19, 45, 0.05)',
  color: '$gray600',
  background: 'rgba(249, 249, 249, 0.75)',
  backdropFilter: 'blur(14px)',
  transition: 'all 0.5s',
  height: '$sizes$header',
})
