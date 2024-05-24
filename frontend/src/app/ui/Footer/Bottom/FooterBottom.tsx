import React from 'react'

import { Flex, Icon, Link, Txt } from '~/shared/ui'

import { StyledFlex, StyledHr } from './FooterBottom.styles'

export const FooterBottom: React.FC = () => {
  return (
    <StyledFlex
      fullWidth
      gap="$3"
      flexWrap="wrap"
      justifyContent="space-between"
    >
      <StyledFlex
        gap="$4"
        alignItems="center"
        flexWrap="wrap"
      >
        <Txt secondary2 color="gray400" css={{ fontFamily: '$body' }}>
          ©
          {' '}
          {new Date().getFullYear()}
          {' '}
          FileMarket Labs Ltd.
        </Txt>
        <StyledHr />
        <Flex gap="$3">
          <Link
            download
            size="small"
            color="gray300"
            href="/docs/PrivacyPolicy.docx"
          >
            Privacy policy
          </Link>
          <Link
            download
            size="small"
            color="gray300"
            href="/docs/TermsOfService.docx"
          >
            Terms of Service
          </Link>
        </Flex>
      </StyledFlex>
      <Flex alignItems="center" gap="$2">
        <Link href="mailto:genesis@filemarket.xyz">
          <Flex gap="$1" alignItems="center">
            <Icon name="email" size={24} />
            <Txt secondary2 color="white" css={{ fontFamily: '$body' }}>
              genesis@filemarket.xyz
            </Txt>
          </Flex>
        </Link>
      </Flex>
    </StyledFlex>
  )
}
