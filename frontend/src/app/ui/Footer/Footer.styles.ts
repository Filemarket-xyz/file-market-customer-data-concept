import { styled } from '~/shared/styles'
import { Container } from '~/shared/ui'

export const StyledFooter = styled('footer', {
  width: '100%',
  backdropFilter: 'blur(12.5px)',
  boxShadow: '0px -4px 15px rgba(19, 19, 45, 0.05)',
  background: '#131416',
  color: '$white',
})

export const StyledContainer = styled(Container, {
  height: '100%',
  display: 'flex',
  justifyContent: 'space-between',
  flexDirection: 'column',
  alignItems: 'center',
  paddingTB: '48px',
  '@md': {
    justifyContent: 'center',
    flexDirection: 'column',
    paddingTop: '32px',
    paddingBottom: '16px',
  },
  '@sm': {
    alignItems: 'center',
  },
})

export const StyledHr = styled('hr', {
  border: 'none',
  width: '100%',
  height: '1px',
  background: '#232528',
  margin: '32px 0',
  '@md': {
    margin: '16px 0',
  },
})
