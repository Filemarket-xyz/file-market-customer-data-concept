import { useQuery } from '@tanstack/react-query'

import { datasetKeys } from '../lib'

export const createMock = () => ({
  dataset: Array.from({ length: 10 }, (_, i) => ({
    name: `Dataset ${i + 1}`,
    data: `mock data ${i + 1}`,
    desc: `mock description ${i + 1}`,
  })),
})

export const useDatasetQuery = () => {
  return useQuery({
    queryKey: datasetKeys.list(),
    // queryFn: () => api.client.datasetGetList().then(({ data }) => data),
    queryFn: () => createMock(),
    select: ({ dataset }) => dataset,
  })
}
