import { serializeError } from '@metamask/rpc-errors'

const fallbackError = { code: 500, message: 'unknown' }

export enum ProviderErrorMessages {
  InternalError = 'Internal JSON-RPC error.',
  InsufficientBalance = 'Actor balance less than needed.',
}

export enum ErrorMessages {
  InsufficientBalance = 'Balance too low for transaction.',
  RejectedByUser = 'Transaction rejected by user.',
}

export const stringifyWeb3Error = (error: any) => {
  if (error?.code === 'ACTION_REJECTED') {
    return ErrorMessages.RejectedByUser
  }

  let message = 'Unknown'
  const serializedError = serializeError(error, { fallbackError })
  const { data }: any = serializedError
  if (serializedError.code === 500) {
    if (data?.cause?.error?.data?.message) {
      const rawMessage: string = data.cause.error.data.message
      // vm error is truncated and useless for us
      message = rawMessage.split(', vm error:')[0]
    } else if (data?.cause?.reason) {
      message = data.cause.reason
    } else if (data?.cause?.shortMessage) {
      message = data.cause.shortMessage
    } else if (data?.cause?.message) {
      message = data.cause.message
    }
  } else if (serializedError.message === ProviderErrorMessages.InternalError) {
    if (data.message?.includes(ProviderErrorMessages.InsufficientBalance)) {
      message = ErrorMessages.InsufficientBalance
    } else {
      message = data.message
    }
  } else {
    message = `Transaction failed. Reason: ${serializedError.message}`
  }

  return `${message} Please try again.`
}
