import 'reset-css'

import { PropsWithChildren } from 'react'

import { globalStyles } from './global'

export const StitchesProvider: React.FC<PropsWithChildren> = ({ children }) => {
  globalStyles()

  return children
}
