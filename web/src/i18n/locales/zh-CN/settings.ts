import type { enUSSettings } from '../en-US/settings'

export const zhCNSettings = {
  'settings.configJson': '配置 JSON',
  'settings.deleteWorkspace': '删除工作区',
  'settings.description': '描述',
  'settings.labelsJson': '标签 JSON',
  'settings.language.description': '该偏好会立即应用到当前工作区界面。',
  'settings.metadata': '元数据',
  'settings.metadata.notLoaded': '工作区元数据尚未加载。',
  'settings.name': '名称',
  'settings.preferences': '偏好设置',
  'settings.replaceConfig': '替换配置',
  'settings.replaceLabels': '替换标签',
  'settings.saved': '已保存',
  'settings.title': '工作区设置',
} satisfies Record<keyof typeof enUSSettings, string>
