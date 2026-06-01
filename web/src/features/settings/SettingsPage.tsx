import { useEffect, useState } from 'react'
import { Save, Trash2 } from 'lucide-react'
import type { WorkspaceMetadata } from '../../api/types'
import { UModelApi } from '../../api/client'
import { Badge, Button, Field, JsonEditor, Panel, TextInput } from '../../design/components'
import { LanguageSelect, useI18n } from '../../i18n'
import { formatError, parseJson, stringify } from '../../lib/json'

type SettingsStatus = { type: 'saved' } | { type: 'error'; message: string }

export function SettingsPage({
  api,
  workspaceId,
  workspace,
  onWorkspaceChange,
  onBack,
}: {
  api: UModelApi
  workspaceId: string
  workspace: WorkspaceMetadata | null
  onWorkspaceChange: (workspace: WorkspaceMetadata | null) => void
  onBack: () => void
}) {
  const { t } = useI18n()
  const [name, setName] = useState(workspace?.name || workspaceId)
  const [description, setDescription] = useState(workspace?.description || '')
  const [labels, setLabels] = useState(stringify(workspace?.labels || {}))
  const [config, setConfig] = useState(stringify(workspace?.config || {}))
  const [replaceLabels, setReplaceLabels] = useState(true)
  const [replaceConfig, setReplaceConfig] = useState(true)
  const [status, setStatus] = useState<SettingsStatus | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    setName(workspace?.name || workspaceId)
    setDescription(workspace?.description || '')
    setLabels(stringify(workspace?.labels || {}))
    setConfig(stringify(workspace?.config || {}))
  }, [workspace, workspaceId])

  async function save() {
    setBusy(true)
    setStatus(null)
    try {
      const next = await api.updateWorkspace(workspaceId, {
        name,
        description,
        labels: parseJson<Record<string, string>>(labels, t('settings.labelsJson')),
        config: parseJson<Record<string, Record<string, unknown>>>(config, t('settings.configJson')),
        if_match_version: workspace?.resource_version,
        replace_labels: replaceLabels,
        replace_config: replaceConfig,
      })
      onWorkspaceChange(next)
      setStatus({ type: 'saved' })
    } catch (error) {
      setStatus({ type: 'error', message: formatError(error) })
    } finally {
      setBusy(false)
    }
  }

  async function remove() {
    setBusy(true)
    setStatus(null)
    try {
      await api.deleteWorkspace(workspaceId)
      onWorkspaceChange(null)
      onBack()
    } catch (error) {
      setStatus({ type: 'error', message: formatError(error) })
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="two-column">
      <Panel
        title={<strong>{t('settings.title')}</strong>}
        action={workspace && <Badge tone={workspace.status === 'active' ? 'success' : 'warning'}>v{workspace.resource_version}</Badge>}
      >
        <div className="stack">
          <Field label={t('settings.name')}>
            <TextInput value={name} onChange={(event) => setName(event.target.value)} />
          </Field>
          <Field label={t('settings.description')}>
            <TextInput value={description} onChange={(event) => setDescription(event.target.value)} />
          </Field>
          <Field label={t('settings.labelsJson')}>
            <JsonEditor value={labels} onChange={setLabels} minHeight={150} />
          </Field>
          <label className="row small muted">
            <input type="checkbox" checked={replaceLabels} onChange={(event) => setReplaceLabels(event.target.checked)} />
            {t('settings.replaceLabels')}
          </label>
          <Field label={t('settings.configJson')}>
            <JsonEditor value={config} onChange={setConfig} minHeight={190} />
          </Field>
          <label className="row small muted">
            <input type="checkbox" checked={replaceConfig} onChange={(event) => setReplaceConfig(event.target.checked)} />
            {t('settings.replaceConfig')}
          </label>
          {status && (
            <Badge tone={status.type === 'saved' ? 'success' : 'danger'}>
              {status.type === 'saved' ? t('settings.saved') : status.message}
            </Badge>
          )}
          <div className="toolbar">
            <Button variant="danger" onClick={() => void remove()} disabled={busy}>
              <Trash2 size={15} />
              {t('settings.deleteWorkspace')}
            </Button>
            <Button variant="primary" onClick={() => void save()} disabled={busy}>
              <Save size={15} />
              {t('common.save')}
            </Button>
          </div>
        </div>
      </Panel>

      <div className="stack">
        <Panel title={<strong>{t('settings.preferences')}</strong>}>
          <div className="stack">
            <LanguageSelect />
            <p className="small muted" style={{ margin: 0 }}>{t('settings.language.description')}</p>
          </div>
        </Panel>

        <Panel title={<strong>{t('settings.metadata')}</strong>}>
          <pre className="result-box small">{workspace ? stringify(workspace) : t('settings.metadata.notLoaded')}</pre>
        </Panel>
      </div>
    </div>
  )
}
