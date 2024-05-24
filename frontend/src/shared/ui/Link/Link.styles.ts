import { styled } from '@stitches/react'
import { ComponentType } from 'react'
import { NavLink } from 'react-router-dom'

import { textVariant } from '../Txt'

export const linkStyled = <Type extends keyof JSX.IntrinsicElements | ComponentType<any>>(element: Type) => {
  return styled(element, {
    fontFamily: '$button',
    fontSize: '$button1',
    fontWeight: '$button',
    lineHeight: '$button',
    color: '$blue500',
    fill: '$blue500',
    transition: 'all 0.25s ease',
    outline: 'none',
    textDecoration: 'none',
    cursor: 'pointer',
    width: 'fit-content',

    '&[data-hovered=true]': {
      opacity: 0.7,
    },
    '&[data-focus=true]': {
      focusRing: '$blue500',
    },
    '&[data-pressed=true]': {
      opacity: 0.9,
    },
    '&[data-disabled=true]': {
      color: '$gray400',
      fill: '$gray400',
      cursor: 'not-allowed',
    },
    variants: {
      underline: {
        true: {
          textDecoration: 'underline',
        },
      },
      color: {
        error: {
          color: '$red',
          fill: '$red',
        },
        gray500: {
          color: '$gray500',
          fill: '$gray500',
        },
        gray300: {
          color: '$gray300',
          fill: '$gray300',
        },
      },
      size: {
        small: {
          ...textVariant('primary2').true,
          fontSize: '14px',
          fontWeight: '400',
        },
      },
    },
  })
}

export const StyledAnhor = linkStyled('a')
export const StyledNavLink = linkStyled(NavLink)
