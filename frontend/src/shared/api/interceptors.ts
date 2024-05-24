import { type MutationCache, type QueryCache } from '@tanstack/react-query'
import { enqueueSnackbar } from 'notistack'
import { NavigateFunction } from 'react-router'

import { AuthResponse, ErrorResponse, HttpResponse } from '~/swagger/Api'

import { RequestMultiplex, stringifyError, stringifyWeb3Error, to } from '../lib'
import { api } from './api'

const refreshMultiplex = new RequestMultiplex<HttpResponse<AuthResponse, ErrorResponse>>()

export const refreshCookies = async (navigate: NavigateFunction): Promise<boolean> => {
  try {
    await refreshMultiplex.request(api.auth.refreshCreate)

    return true
  } catch (error) {
    navigate(to.root())

    return false
  }
}

const showError = (response: HttpResponse<unknown, ErrorResponse>) => {
  const message = 'error' in response ? stringifyError(response.error) : stringifyWeb3Error(response)

  return enqueueSnackbar(message, {
    variant: 'error',
  })
}

export const onQueryError: (navigate: NavigateFunction) => QueryCache['config']['onError'] = (navigate) => {
  return async (response, query) => {
    if (response.status == 401) {
      const isRefreshSuccess = await refreshCookies(navigate)

      if (isRefreshSuccess) {
        return query.fetch(query.options)
      }
    }

    showError(response)
  }
}

export const onMutationError: (navigate: NavigateFunction) => MutationCache['config']['onError'] = (navigate) => {
  return async (response, variables, _, mutation) => {
    if (response.status === 401) {
      const isRefreshSuccess = await refreshCookies(navigate)

      if (isRefreshSuccess) {
        return mutation.execute(variables)
      }
    }

    showError(response)
  }
}
