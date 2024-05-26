import { Logo } from '~/entities/Logo'
import { to } from '~/shared/lib'
import { Flex, Icon, IconName, Link, LinkProps, Txt } from '~/shared/ui'

import { StyledFirstColumnFlex, StyledH4, StyledLastColumnFlex, StyledLink, StyledMiddleColumnFlex, StyledWrapperFlex } from './FooterTop.styles'

interface FooterLink extends LinkProps {
  iconName: IconName
}

interface FooterLinkColumn {
  title: string
  links: LinkProps[]
}

const buttons: FooterLink[] = [
  {
    iconName: 'twitter',
    children: 'Twitter',
    href: 'https://twitter.com/filemarket_xyz',
  },
  {
    iconName: 'discord',
    children: 'Discord',
    href: 'https://discord.gg/filemarket',
  },
  {
    iconName: 'telegram',
    children: 'Telegram',
    href: 'https://t.me/FileMarketChat',
  },
  {
    iconName: 'youtube',
    children: 'Youtube',
    href: 'https://www.youtube.com/@filemarket_xyz',
  },
  {
    iconName: 'medium',
    children: 'Medium',
    href: 'https://medium.com/filemarket-xyz',
  },
  {
    iconName: 'linkedin',
    children: 'LinkedIn',
    href: 'https://www.linkedin.com/company/filemarketxyz/',
  },
]

const columns: FooterLinkColumn[] = [
  {
    title: 'Platform',
    links: [
      {
        children: 'Explore EFTs',
        href: 'https://filemarket.xyz/market',
      },
      {
        children: 'Build own Shop',
        href: 'https://form.typeform.com/to/gulmhUKG?typeform-source=filemarket.xyz',
      },
      {
        children: 'Collections',
        href: 'https://filemarket.xyz/market/collections',
      },
      {
        children: 'FileBunnies',
        href: 'https://filemarket.xyz/fileBunnies',
      },
      {
        children: 'How to get FIL',
        href: 'https://medium.com/filemarket-xyz/how-to-buy-fil-and-use-fil-in-the-filecoin-virtual-machine-d67fa90764d5',
      },
      {
        children: 'FAQ',
        to: to.root(),
      },
    ],
  },
  {
    title: 'Links',
    links: [
      {
        children: 'EFT Protocol',
        href: 'https://medium.com/filemarket-xyz/how-to-attach-an-encrypted-file-to-your-nft-7d6232fd6d34',
      },
      {
        children: 'SDK',
        to: to.root(),
      },
      {
        children: 'DAO',
        href: 'https://discord.gg/filemarket',
      },
      {
        children: 'GitHub',
        href: 'https://github.com/Filemarket-xyz/file-market',
      },
      {
        children: 'Blogs',
        href: 'https://filemarket.ghost.io/',
      },
    ],
  },
  {
    title: 'Company',
    links: [
      {
        children: 'About',
        to: to.root(),
      },
      {
        children: 'Ambassador program',
        href: 'https://filemarket.typeform.com/to/MTwDOB1J',
      },
      {
        children: 'Become a partner',
        href: 'https://filemarket.typeform.com/to/BqkdzJQM',
      },
      {
        children: 'Branding',
        href: 'https://filemarket.xyz/branding',
      },
      {
        children: 'Calendly',
        href: 'http://calendly.com/filemarket',
      },
    ],
  },
]

export const FooterTop: React.FC = () => {
  return (
    <StyledWrapperFlex
      fullWidth
      flexWrap="wrap"
    >
      <StyledFirstColumnFlex
        gap={16}
        flexDirection="column"
      >
        <Logo color="$white" />
        <Txt secondary1 color="white" css={{ lineHeight: 1.5 }}>
          FileMarket is a multi-chain marketplace specializing in the
          tokenization and monetization of pivotal public data through perpetual
          decentralized storage with a privacy layer, opening the data economy
          to the mass market.
        </Txt>
      </StyledFirstColumnFlex>
      {columns.map(({ title, links }) => (
        <StyledMiddleColumnFlex
          key={title}
          gap={16}
          flexDirection="column"
        >
          <StyledH4>{title}</StyledH4>
          <Flex flexDirection="column" gap="6px">
            {links.map(({ children, href, ...props }, i) => (
              <Link
                key={i}
                {...props}
                href={href}
                target={href ? '_blank' : '_self'}
              >
                <Txt secondary2 color="white" css={{ fontWeight: 500 }}>
                  {children}
                </Txt>
              </Link>
            ))}
          </Flex>
        </StyledMiddleColumnFlex>
      ))}
      <StyledLastColumnFlex
        gap={16}
        flexDirection="column"
      >
        <StyledH4>Join our community</StyledH4>
        <Flex fullWidth gap="4px" flexWrap="wrap">
          {buttons.map(({ iconName, href, children }, index) => (
            <StyledLink key={index} href={href} target="_blank">
              <Icon name={iconName} />
              <Txt secondary2 color="white" css={{ fontWeight: 500 }}>
                {children}
              </Txt>
            </StyledLink>
          ))}
        </Flex>
      </StyledLastColumnFlex>
    </StyledWrapperFlex>
  )
}
