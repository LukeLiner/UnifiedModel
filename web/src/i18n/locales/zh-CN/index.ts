import type { enUS } from '../en-US'
import { zhCNCommon } from './common'
import { zhCNLanding } from './landing'
import { zhCNSettings } from './settings'

export const zhCN = {
  ...zhCNCommon,
  ...zhCNSettings,
  ...zhCNLanding,
} satisfies Record<keyof typeof enUS, string>
