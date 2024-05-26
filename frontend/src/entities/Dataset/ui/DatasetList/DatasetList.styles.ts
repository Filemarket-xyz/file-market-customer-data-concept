import { AppCSS, styled } from '~/shared/styles'

import { rowCss } from './Item'

export const StyledRoot = styled('div', {
  width: '100%',
  overflowX: 'auto',
})

export const StyledHeadline = styled('div', rowCss, {
  color: '$gray400',
  marginBottom: 16,
})

export const listCss: AppCSS = {
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',
  width: '100%',
}
