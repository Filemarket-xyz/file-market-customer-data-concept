import React from 'react'

import { Card, Txt } from '~/shared/ui'
import { Dataset } from '~/swagger/Api'

import { StyledRow } from './DatasetItem.styles'

export interface DatasetItemProps {
  dataset?: Dataset
}

export const DatasetItem: React.FC<DatasetItemProps> = ({ dataset }) => {
  return (
    <Card disableShadow padding="24px 16px">
      <StyledRow>
        <Txt primary1>{dataset?.name}</Txt>
        <Txt primary1 color="gray800">
          &quot;
          {dataset?.data}
          &quot;
        </Txt>
        <Txt primary1>{dataset?.desc}</Txt>
      </StyledRow>
    </Card>
  )
}
