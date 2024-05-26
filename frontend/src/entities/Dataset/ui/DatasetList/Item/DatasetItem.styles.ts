import { css, styled } from '~/shared/styles'

export const rowCss = css({
  display: 'grid',
  gap: 12,
  gridTemplateColumns: '1fr 1fr 2fr',
})

export const StyledRow = styled('div', rowCss, {
  color: '$gray500',
})
