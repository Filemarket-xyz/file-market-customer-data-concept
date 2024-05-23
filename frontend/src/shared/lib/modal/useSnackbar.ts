import { useSnackbar as useNotistackSnackbar } from 'notistack'
import { ReactNode, useMemo } from 'react'

export const useSnackbar = () => {
  const { enqueueSnackbar } = useNotistackSnackbar()

  return useMemo(() => ({
    error: (body: ReactNode) => enqueueSnackbar(body, {
      autoHideDuration: 2000,
      variant: 'error',
    }),
    success: (body: ReactNode) => enqueueSnackbar(body, {
      autoHideDuration: 1500,
      variant: 'success',
    }),
  }), [enqueueSnackbar])
}
