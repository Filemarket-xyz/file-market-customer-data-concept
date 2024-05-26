import { keyframes, styled } from '~/shared/styles'

const spin = keyframes({
  '0%': {
    transform: 'translate(-50%, -50%) rotate(0deg)',
  },
  '100%': {
    transform: 'translate(-50%, -50%) rotate(360deg)',
  },
})

export const StyledLayout = styled('div', {
  background: 'transparent',
})

export const StyledRoot = styled('div', {
  borderRadius: 16,
  border: '2px solid #6B6F76',
  variants: {
    variant: {
      elevate: {
        background: '$white',
        boxShadow: '8px 8px 0px 0px #13141540',
      },
      outline: {
        borderColor: '$gray200',
        background: 'gray100',
        boxShadow: '0px 0px 15px 0px #13132D0D',
      },
      gradient: {
        overflow: 'hidden',
        position: 'relative',
        border: 'none',

        '&::before': {
          content: '""',
          position: 'absolute',
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -50%)',
          width: 1500,
          height: 1500,
          transformOrigin: 'center center',
          backgroundImage: 'linear-gradient(180deg, #8efdb5 30%, #028fff 50%, #01e3f8 70%)',
          animation: `${spin} 2.5s linear infinite`,
        },

        [`& > ${StyledLayout.selector}`]: {
          position: 'relative',
          borderRadius: 14,
          zIndex: 1,
          margin: 2,
          background: '$gray100',
        },
      },
    },
  },
})
