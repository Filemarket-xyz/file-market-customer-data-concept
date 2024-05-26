import { styled } from '~/shared/styles'

import Discord from './discord.svg?react'
import Email from './email.svg?react'
import Linkedin from './linkedin.svg?react'
import Logo from './logo.svg?react'
import Medium from './medium.svg?react'
import Storage from './storage.svg?react'
import Telegram from './telegram.svg?react'
import Twitter from './twitter.svg?react'
import Youtube from './youtube.svg?react'

export const icons = {
  email: { Component: styled(Email) },
  discord: { Component: styled(Discord) },
  linkedin: { Component: styled(Linkedin) },
  medium: { Component: styled(Medium) },
  telegram: { Component: styled(Telegram) },
  twitter: { Component: styled(Twitter) },
  youtube: { Component: styled(Youtube) },
  logo: { Component: styled(Logo) },
  storage: { Component: styled(Storage) },
} as const

export type IconName = keyof typeof icons
