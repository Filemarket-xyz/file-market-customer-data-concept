import React, { ComponentType } from 'react'

export const withSuspense = <T extends JSX.IntrinsicAttributes >(WrappedComponent: ComponentType<T>) => {
  return (props: T) => (
    <React.Suspense fallback="loading...">
      <WrappedComponent {...props} />
    </React.Suspense>
  )
}
