import React from 'react'

import { Txt, VirtualScroll } from '~/shared/ui'

import { useDatasetQuery } from '../../model'
import { listCss, StyledHeadline, StyledRoot } from './DatasetList.styles'
import { DatasetItem } from './Item'

export const DatasetList: React.FC = () => {
  const datasetQuery = useDatasetQuery()

  return (
    <StyledRoot>
      <StyledHeadline>
        <Txt primary2>Field name</Txt>
        <Txt primary2>Value</Txt>
        <Txt primary2>Description</Txt>
      </StyledHeadline>
      <VirtualScroll
        hasMore={false}
        loading={datasetQuery.isLoading || datasetQuery.isFetching}
        currentItemCount={datasetQuery.data?.length ?? 0}
        fetchMore={console.log}
        render={({ index }) => <DatasetItem dataset={datasetQuery.data?.[index]} />}
        listCss={listCss}
      />
    </StyledRoot>
  )
}
