import { globalCss } from './stitches.config'

export const globalStyles = globalCss({
  '@font-face': [
    {
      fontFamily: '',
      fontStyle: 'normal',
      fontWeight: 400,
      src: `local(''), url('') format('truetype')`,
    },
  ],
  'html, body, #root': {
    height: '100svh',
    minHeight: '100svh',
  },
  a: {
    textDecoration: 'none',
    color: 'inherit',
  },
})
