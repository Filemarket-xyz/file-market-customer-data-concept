import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React, { useMemo } from 'react'
import { Outlet, useNavigate } from 'react-router'

import { onMutationError, onQueryError } from '~/shared/api'
import { ModalProvider } from '~/shared/lib'
import { StitchesProvider } from '~/shared/styles'

export const ProvidersOutlet: React.FC = () => {
  const navigate = useNavigate()

  const queryClient = useMemo(() => new QueryClient({
    queryCache: new QueryCache({
      onError: onQueryError(navigate),
    }),
    mutationCache: new MutationCache({
      onError: onMutationError(navigate),
    }),
  }), [navigate])

  return (
    <StitchesProvider>
      <QueryClientProvider client={queryClient}>
        <ModalProvider>
          <Outlet />
        </ModalProvider>
      </QueryClientProvider>
    </StitchesProvider>
  )
}
