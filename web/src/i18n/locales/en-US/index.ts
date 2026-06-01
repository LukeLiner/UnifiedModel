import { enUSCommon } from './common'
import { enUSLanding } from './landing'
import { enUSSettings } from './settings'

export const enUS = {
  ...enUSCommon,
  ...enUSSettings,
  ...enUSLanding,
} as const
