import React from 'react'
import { Outlet } from 'react-router-dom'

import { StyledBody, StyledLayout } from './Layout.styles'

export const Layout: React.FC = () => {
  return (
    <StyledLayout>
      <StyledBody>
        <Outlet />
      </StyledBody>
    </StyledLayout>
  )
}
