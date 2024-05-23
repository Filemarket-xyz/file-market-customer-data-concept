import { createBrowserRouter, Navigate } from "react-router-dom"

import { Layout } from '~/app/ui'
import { Routes, to } from '~/shared/lib'

import { ProvidersOutlet } from '../ProvidersOutlet'

export const router = createBrowserRouter([
  {
    path: Routes.Root,
    element: <ProvidersOutlet />,
    children: [
      {
        path: Routes.Root,
        element: <Layout />,
        children: [
          { path: Routes.Root, element: <>root page</> },
        ],
      },

      { path: '*', element: <Navigate replace to={to.root()} /> },
    ],
  },
])
