import React from 'react'
import { Outlet } from 'react-router-dom'

import { Footer } from '../Footer'
import { Header } from '../Header'
import { StyledBody, StyledLayout } from './Layout.styles'

export const Layout: React.FC = () => {
  return (
    <StyledLayout>
      <Header />
      <StyledBody>
        <Outlet />
      </StyledBody>
      <Footer />
    </StyledLayout>
  )
}
