import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { SnackbarProvider } from 'notistack'
import React, { useMemo } from 'react'
import { Outlet, useNavigate } from 'react-router'
import { ThirdwebProvider } from 'thirdweb/react'
import { createConfig, http, WagmiProvider } from 'wagmi'
import { mainnet } from 'wagmi/chains'

import { onMutationError, onQueryError } from '~/shared/api'
import { ModalProvider } from '~/shared/lib'
import { StitchesProvider } from '~/shared/styles'

const config = createConfig({
  chains: [mainnet],
  transports: {
    [mainnet.id]: http(),
  },
})

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
    <ThirdwebProvider>
      <SnackbarProvider>
        <StitchesProvider>
          <WagmiProvider config={config}>
            <QueryClientProvider client={queryClient}>
              <ModalProvider>
                <Outlet />
              </ModalProvider>
            </QueryClientProvider>
          </WagmiProvider>
        </StitchesProvider>
      </SnackbarProvider>
    </ThirdwebProvider>
  )
}
