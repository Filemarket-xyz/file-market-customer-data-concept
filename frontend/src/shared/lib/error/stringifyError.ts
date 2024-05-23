import { ErrorResponse } from '~/swagger/Api'

import { errorResponseToMessage } from './errorResponseToMessage'

/**
 * Safely converts any error to string
 * @param error an error
 */
export function stringifyError(error: unknown): string {
  if (typeof error === 'string') return error

  if (
    error instanceof Object &&
    'message' in error &&
    typeof error.message === 'string' &&
    'detail' in error &&
    typeof error.detail === 'string'
  ) {
    // Error response
    return errorResponseToMessage(error as ErrorResponse)
  }

  if (error instanceof Error) {
    if (error.stack) return error.stack

    return `${error.name}: ${error.message}`
  }

  let message
  try {
    message = JSON.stringify(error)
    if (message === '{}') {
      message = `${error}`
    }
  } catch (e) {
    message = `${error}`
  }

  return message
}
