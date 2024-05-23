import { ErrorResponse } from '~/swagger/Api'

export function errorResponseToMessage(error?: ErrorResponse): string {
  if (!error) {
    return 'received nullish error from the backend, but request was not successful'
  }

  return `${error.message}${error.detail ? `: ${error.detail}` : ''}`
}
