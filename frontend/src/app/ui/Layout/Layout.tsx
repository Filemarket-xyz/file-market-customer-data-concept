import React from 'react'
import { Outlet } from 'react-router-dom'

import { Header } from '../Header'
import { StyledBody, StyledLayout } from './Layout.styles'

export const Layout: React.FC = () => {
  return (
    <StyledLayout>
      <Header />
      <StyledBody>
        <Outlet />
      </StyledBody>
    </StyledLayout>
  )
}
