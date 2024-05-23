import { styled } from '~/shared/styles'

export const StyledDiv = styled('div', {
  display: 'flex',
  variants: {
    fullWidth: {
      true: {
        width: '100%',
      },
    },
    fullHeight: {
      true: {
        height: '100%',
      },
    },
  },
})
