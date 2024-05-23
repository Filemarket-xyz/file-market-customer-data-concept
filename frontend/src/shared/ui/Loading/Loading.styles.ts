import { keyframes, styled } from '~/shared/styles'

const spin = keyframes({
  '0%': {
    rotate: '0deg',
  },
  '100%': {
    rotate: '360deg',
  },
})

export const StyledDiv = styled('div', {
  border: '3px solid $gray300',
  borderTop: '3px solid $blue500',
  borderRadius: '50%',
  animation: `${spin} 0.6s linear infinite`,
  variants: {
    color: {
      primary: {
        borderTopColor: '$blue500',
      },
      disabled: {
        borderColor: '$gray400',
        borderTopColor: '$gray500',
      },
    },
  },
})
