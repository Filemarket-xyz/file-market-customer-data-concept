/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

import { ErrorResponse, HttpResponse } from './swagger/Api'

declare module '@tanstack/react-query' {
  interface Register {
    defaultError: HttpResponse<unknown, ErrorResponse>
  }
}
