import { reset } from 'stitches-reset'

import InterMediumTtf from '../fonts/Inter-Medium.ttf'
import InterMediumWoff from '../fonts/Inter-Medium.woff'
import InterMediumWoff2 from '../fonts/Inter-Medium.woff2'
import InterRegularTtf from '../fonts/Inter-Regular.ttf'
import InterRegularWoff from '../fonts/Inter-Regular.woff'
import InterRegularWoff2 from '../fonts/Inter-Regular.woff2'
import MontserratBoldTtf from '../fonts/Montserrat-Bold.ttf'
import MontserratBoldWoff from '../fonts/Montserrat-Bold.woff'
import MontserratBoldWoff2 from '../fonts/Montserrat-Bold.woff2'
import MontserratRegularTtf from '../fonts/Montserrat-Regular.ttf'
import MontserratRegularWoff from '../fonts/Montserrat-Regular.woff'
import MontserratRegularWoff2 from '../fonts/Montserrat-Regular.woff2'
import MontserratSemiBoldTtf from '../fonts/Montserrat-SemiBold.ttf'
import MontserratSemiBoldWoff from '../fonts/Montserrat-SemiBold.woff'
import MontserratSemiBoldWoff2 from '../fonts/Montserrat-SemiBold.woff2'
import MuseoModernoBoldTtf from '../fonts/MuseoModerno-Bold.ttf'
import MuseoModernoBoldWoff from '../fonts/MuseoModerno-Bold.woff'
import MuseoModernoBoldWoff2 from '../fonts/MuseoModerno-Bold.woff2'
import MuseoModernoSemiBoldTtf from '../fonts/MuseoModerno-SemiBold.ttf'
import MuseoModernoSemiBoldWoff from '../fonts/MuseoModerno-SemiBold.woff'
import MuseoModernoSemiBoldWoff2 from '../fonts/MuseoModerno-SemiBold.woff2'
import { globalCss } from './stitches.config'

export const globalStyles = globalCss({
  ...reset,
  '@font-face': [
    {
      fontFamily: 'Inter',
      fontStyle: 'normal',
      fontWeight: 400,
      src: `url('${InterRegularWoff2}') format('woff2'),
            url('${InterRegularWoff}') format('woff'),
            url('${InterRegularTtf}') format('truetype')`,
      fontDisplay: 'swap',
    },
    {
      fontFamily: 'Inter',
      fontStyle: 'normal',
      fontWeight: 500,
      src: `url('${InterMediumWoff2}') format('woff2'),
            url('${InterMediumWoff}') format('woff'),
            url('${InterMediumTtf}') format('truetype')`,
      fontDisplay: 'swap',
    },
    {
      fontFamily: 'Montserrat',
      fontStyle: 'normal',
      fontWeight: 400,
      src: `url('${MontserratRegularWoff2}') format('woff2'),
            url('${MontserratRegularWoff}') format('woff'),
            url('${MontserratRegularTtf}') format('truetype')`,
    },
    {
      fontFamily: 'Montserrat',
      fontStyle: 'normal',
      fontWeight: 500,
      src: `url('${MontserratSemiBoldWoff2}') format('woff2'),
            url('${MontserratSemiBoldWoff}') format('woff'),
            url('${MontserratSemiBoldTtf}') format('truetype')`,
    },
    {
      fontFamily: 'Montserrat',
      fontStyle: 'normal',
      fontWeight: 600,
      src: `url('${MontserratBoldWoff2}') format('woff2'),
            url('${MontserratBoldWoff}') format('woff'),
            url('${MontserratBoldTtf}') format('truetype')`,
    },
    {
      fontFamily: 'MuseoModerno',
      fontStyle: 'normal',
      fontWeight: 600,
      src: `url('${MuseoModernoSemiBoldWoff2}') format('woff2'),
            url('${MuseoModernoSemiBoldWoff}') format('woff'),
            url('${MuseoModernoSemiBoldTtf}') format('truetype')`,
    },
    {
      fontFamily: 'MuseoModerno',
      fontStyle: 'normal',
      fontWeight: 700,
      src: `url('${MuseoModernoBoldWoff2}') format('woff2'),
            url('${MuseoModernoBoldWoff}') format('woff'),
            url('${MuseoModernoBoldTtf}') format('truetype')`,
    },
  ],
  'html, body, #root': {
    fontFamily: '$primary',
    fontSize: '$html',
    height: '100svh',
    minHeight: '100svh',
  },
  '*::-webkit-scrollbar, html *::-webkit-scrollbar': {
    width: '10px',
    height: '4px',
  },
  '*::-webkit-scrollbar-track, html *::-webkit-scrollbar-track': {
    background: 'none',
    boxShadow: 'inset 0 0 5px 5px #0090FF',
    border: 'solid 6px transparent',
  },
  '*::-webkit-scrollbar-thumb, html *::-webkit-scrollbar-thumb': {
    background: '#0090FF',
    borderRadius: '8px',
    border: 'solid 1px rgba(255, 255, 255, 0.5)',
  },
  a: {
    textDecoration: 'none',
    color: 'inherit',
  },
  '*, *:before, *:after': {
    boxSizing: 'border-box',
  },
})
