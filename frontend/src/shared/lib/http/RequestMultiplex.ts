interface Call<RequestResult> {
  onFulfilled: (value: RequestResult) => void,
  onRejected: (reason: unknown) => void,
}

export class RequestMultiplex<RequestResult> {
  private waitingCalls: Call<RequestResult>[] = []

  private currentRequest: Promise<void> | undefined

  private reset(): void {
    this.waitingCalls = []
    this.currentRequest = undefined
  }

  request(requester: () => Promise<RequestResult>): Promise<RequestResult> {
    return new Promise((resolve, reject) => {
      this.waitingCalls.push({
        onFulfilled: resolve,
        onRejected: reject,
      })

      if (!this.currentRequest) {
        this.currentRequest = requester()
          .then((result) => this.waitingCalls.forEach((call) => call.onFulfilled(result)))
          .catch((error) => this.waitingCalls.forEach((call) => call.onRejected(error)))
          .finally(() => this.reset())
      }
    })
  }
}
