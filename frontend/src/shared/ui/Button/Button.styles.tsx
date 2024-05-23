import { styled } from '~/shared/styles'

import { textVariant } from '../Txt'

export const StyledButton = styled('button', {
  ...textVariant('button1').true,
  height: '48px',
  minWidth: '160px',
  outline: 'none',
  border: 'none',
  borderRadius: '$1',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: '0 24px',
  userSelect: 'none',
  position: 'relative',
  overflow: 'hidden',
  cursor: 'pointer',
  transition: 'transform 0.25s ease, opacity 0.25s ease, scale 0.25s ease',
  textDecoration: 'none',
  '&[data-pressed=true]': {
    scale: 0.97,
  },
  '&[data-hovered=true]': {
    opacity: 0.7,
  },
  '&[data-disabled=true]': {
    cursor: 'not-allowed',
  },
  '&[data-focus-ring=true]': {
    focusRing: '$blue500',
  },
  variants: {
    variant: {
      primary: {
        color: '$white',
        fill: '$white',
        background: '$gradients$mainNew',
        '&[data-disabled=true]': {
          background: '$gray400',
          color: '$white',
          fill: '$white',
        },
      },
      secondary: {
        color: '$white',
        fill: '$white',
        background: 'linear-gradient(to top, #028FFF 50%, $white 50%)',
        backgroundSize: '100% 200%',
        backgroundPosition: '0 100%',
        transition: 'background-position 0.3s ease-out, color 0.3s ease-out, transform 0.25s ease, opacity 0.25s ease, scale 0.25s ease',
        border: '1px solid #028FFF',
        '&[data-hovered=true]': {
          opacity: 1,
          color: '#028FFF',
          backgroundPosition: '0 0',
        },
        '&[data-focus-ring=true]': {
          focusRing: '$blue500',
        },
        '&[data-disabled=true]': {
          background: '$gray200',
          color: '$gray600',
          fill: '$gray600',
          border: '2px solid $gray300',
        },
      },
      text: {
        color: '$blue500',
        fill: '$blue500',
        background: 'transparent',
        '&[data-disabled=true]': {
          color: '$gray600',
          fill: '$gray600',
        },
      },
      error: {
        border: '2px solid #C54B5C',
        color: '#C54B5C',
        background: 'white',
      },
      glowing: {
        minWidth: 240,
        transition: 'all 0.25s ease',
        color: '$blue500',
        border: '2px solid $blue500',
        background: '$white',
        borderRadius: '$2',
        '&[data-hovered=true]': {
          opacity: 1,
          color: '$blue300',
          borderColor: '$blue300',
        },
      },
    },
    size: {
      small: {
        height: '36px',
        minWidth: 0,
        padding: '0 18px',
      },
    },
    loading: {
      true: {
        gap: '$2',
      },
    },
    fullWidth: {
      true: {
        width: '100%',
      },
    },
    icon: {
      true: {
        color: '$white',
        fill: '$white',
        minWidth: 0,
        padding: 0,
        size: '48px',
        borderRadius: '50%',
        background: 'transparent',
        '& > *': {
          height: '26px',
        },
        '&[data-disabled=true]': {
          color: '$gray600',
          fill: '$gray600',
        },
      },
    },
    isDisabled: {
      true: {
        color: '$gray400',
        border: '2px solid $gray400 !important',
        cursor: 'not-allowed',
        '& img': {
          filter: 'contrast(0), brightness(1.2)',
        },
      },
    },
  },
  compoundVariants: [
    {
      icon: true,
      size: 'small',
      css: {
        size: '36px',
        padding: 0,
        '& > *': {
          height: '20px',
        },
      },
    },
    {
      icon: true,
      variant: 'primary',
      css: {
        color: '$white',
        fill: '$white',
        background: '$gradients$mainNew',
      },
    },
    {
      icon: true,
      variant: 'secondary',
      css: {
        color: '$blue500',
        fill: '$blue500',
        backgroundColor: '$gray100',
        '&[data-disabled=true]': {
          background: '$gray200',
          color: '$gray600',
          fill: '$gray600',
        },
      },
    },
    {
      icon: true,
      variant: 'text',
      css: {
        color: '$blue500',
        fill: '$blue500',
        background: 'transparent',
        '&[data-disabled=true]': {
          color: '$gray600',
          fill: '$gray600',
        },
      },
    },
  ],
})
