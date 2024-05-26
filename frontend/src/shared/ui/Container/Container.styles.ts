import { styled } from '~/shared/styles'

export const Container = styled('div', {
  width: '100%',
  paddingLR: 'calc((100% - $breakpoints$xl) * 0.3 + 32px)',
  '@xl': {
    paddingLR: 'calc((100% - $breakpoints$lg) * 0.554 + 32px)',
  },
  '@lg': {
    paddingLR: 'calc((100% - $breakpoints$md) * 0.554 + 32px)',
  },
  '@md': {
    paddingLR: 'calc((100% - $breakpoints$sm) * 0.554 + 16px)',
  },
  '@sm': {
    paddingLR: 16,
  },
  variants: {
    fullHeight: {
      true: {
        height: '100%',
      },
    },
  },
})
