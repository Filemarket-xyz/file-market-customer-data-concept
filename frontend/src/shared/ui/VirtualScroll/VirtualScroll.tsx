import React, { Fragment, ReactNode, useCallback, useEffect, useMemo, useRef } from 'react'
import useVirtual, { Item } from 'react-cool-virtual'

import { AppCSS } from '~/shared/styles'

import { Loading } from '../Loading'
import { StyledInnerDiv, StyledTrigger } from './VirtualScroll.styles'

export interface VirtualScrollProps {
  currentItemCount: number
  render: (item: Item) => ReactNode
  fetchMore: () => void
  hasMore: boolean
  loading: boolean
  listCss?: AppCSS
}

export const VirtualScroll: React.FC<VirtualScrollProps> = ({
  currentItemCount,
  render,
  hasMore,
  loading,
  fetchMore,
  listCss,
}) => {
  const triggerRef = useRef<HTMLSpanElement>(null)
  const { outerRef, innerRef, items } = useVirtual<HTMLDivElement, HTMLDivElement>({
    itemCount: currentItemCount,
  })

  const observerCallback = useCallback<IntersectionObserverCallback>(([target]) => {
    if (!target.isIntersecting || loading) return

    fetchMore()
  }, [loading])

  const observer = useMemo(() => new IntersectionObserver(observerCallback, { threshold: 0 }), [observerCallback])

  useEffect(() => {
    if (!triggerRef.current) return

    observer.observe(triggerRef.current)

    return () => { observer.disconnect() }
  }, [triggerRef.current, observer])

  return (
    <div ref={outerRef}>
      <StyledInnerDiv
        ref={innerRef}
        css={listCss}
      >
        {/* items.length may be incorrect at first render if there was a larger value at the previous render */}
        {items.length <= currentItemCount && items.map((item) => (
          <Fragment key={item.index}>
            {render(item)}
          </Fragment>
        ))}
        {!!items.length && hasMore && (
          <Loading loading={loading}>
            <StyledTrigger ref={triggerRef} />
          </Loading>
        )}
      </StyledInnerDiv>
    </div>
  )
}
