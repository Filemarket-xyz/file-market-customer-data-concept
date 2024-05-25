import { useQuery } from '@tanstack/react-query'

import { api } from '~/shared/api'
import { AuthResponse, ErrorResponse, UserResponse } from '~/swagger/Api'

import { userKeys } from '../lib'

export const useUserQuery = () => {
  return useQuery<AuthResponse | null, ErrorResponse, UserResponse | undefined | null>({
    queryKey: userKeys.root,
    queryFn: () => api.auth.postAuth().then(({ data }) => data),
    select: (data) => data?.user,
    retry: 1,
  })
}
