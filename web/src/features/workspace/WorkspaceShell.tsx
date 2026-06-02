import { Suspense, lazy, useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import { ArrowLeft, PanelLeftClose, PanelLeftOpen, RefreshCcw } from 'lucide-react'
import type { HealthResponse, WorkspaceMetadata } from '../../api/types'
import { UModelApi } from '../../api/client'
import { Brand, HealthBadge } from '../../App'
import { Button, Badge, IconButton } from '../../design/components'
import { useI18n, type MessageKey } from '../../i18n'
import { formatError } from '../../lib/json'

const ExplorerPage = lazy(() => import('../explorer/ExplorerPage').then(({ ExplorerPage }) => ({ default: ExplorerPage })))
const EntityTopoPage = lazy(() => import('../entityTopo/EntityTopoPage').then(({ EntityTopoPage }) => ({ default: EntityTopoPage })))
const QueryPage = lazy(() => import('../query/QueryPage').then(({ QueryPage }) => ({ default: QueryPage })))
const ImportsPage = lazy(() => import('../imports/ImportsPage').then(({ ImportsPage }) => ({ default: ImportsPage })))
const AgentPage = lazy(() => import('../agent/AgentPage').then(({ AgentPage }) => ({ default: AgentPage })))
const SettingsPage = lazy(() => import('../settings/SettingsPage').then(({ SettingsPage }) => ({ default: SettingsPage })))
const ApiMapPage = lazy(() => import('../settings/ApiMapPage').then(({ ApiMapPage }) => ({ default: ApiMapPage })))
const DataStorePage = lazy(() => import('../query/DataStorePage').then(({ DataStorePage }) => ({ default: DataStorePage })))

export type WorkspaceView = 'explorer' | 'entityTopo' | 'query' | 'imports' | 'agent' | 'settings' | 'docs' | 'data'

interface NavItem {
  value: WorkspaceView
  label: string
  icon: ReactNode
}

export function WorkspaceShell({
  api,
  workspaceId,
  workspace,
  health,
  view,
  navItems,
  onViewChange,
  onWorkspaceChange,
  onHealthChange,
  onBack,
}: {
  api: UModelApi
  workspaceId: string
  workspace: WorkspaceMetadata | null
  health: HealthResponse | null
  view: WorkspaceView
  navItems: NavItem[]
  onViewChange: (view: WorkspaceView) => void
  onWorkspaceChange: (workspace: WorkspaceMetadata | null) => void
  onHealthChange: (health: HealthResponse | null) => void
  onBack: () => void
}) {
  const { t } = useI18n()
  const [error, setError] = useState('')
  const [refreshToken, setRefreshToken] = useState(0)
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false)

  const refresh = useCallback(async () => {
    setError('')
    try {
      const [nextHealth, nextWorkspace] = await Promise.all([
        api.health().catch(() => null),
        api.getWorkspace(workspaceId),
      ])
      onHealthChange(nextHealth)
      onWorkspaceChange(nextWorkspace)
      setRefreshToken((value) => value + 1)
    } catch (nextError) {
      setError(formatError(nextError))
    }
  }, [api, onHealthChange, onWorkspaceChange, workspaceId])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const page = useMemo(() => {
    switch (view) {
      case 'explorer':
        return <ExplorerPage api={api} workspaceId={workspaceId} refreshToken={refreshToken} />
      case 'entityTopo':
        return <EntityTopoPage api={api} workspaceId={workspaceId} refreshToken={refreshToken} />
      case 'query':
        return <QueryPage api={api} workspaceId={workspaceId} />
      case 'imports':
        return <ImportsPage api={api} workspaceId={workspaceId} onChanged={() => setRefreshToken((value) => value + 1)} />
      case 'agent':
        return <AgentPage api={api} workspaceId={workspaceId} />
      case 'settings':
        return (
          <SettingsPage
            api={api}
            workspaceId={workspaceId}
            workspace={workspace}
            onWorkspaceChange={onWorkspaceChange}
            onBack={onBack}
          />
        )
      case 'docs':
        return <ApiMapPage />
      case 'data':
        return <DataStorePage api={api} workspaceId={workspaceId} />
      default:
        return null
    }
  }, [api, onBack, onWorkspaceChange, refreshToken, view, workspace, workspaceId])

  const explorerHost = view === 'explorer' || view === 'entityTopo'
  const topbarHidden = explorerHost || view === 'query' || view === 'imports' || view === 'settings'

  return (
    <div className={`workspace-shell app-shell ${sidebarCollapsed ? 'collapsed' : ''} ${explorerHost ? 'explorer-host' : ''}`}>
      <aside className="workspace-sidebar">
        <div className="workspace-sidebar-header">
          <Brand compact />
          <div className="workspace-sidebar-title" style={{ minWidth: 0 }}>
            <strong style={{ display: 'block', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {workspace?.name || workspaceId}
            </strong>
            <span className="workspace-id">{workspaceId}</span>
          </div>
          <IconButton
            className="workspace-collapse-button"
            label={sidebarCollapsed ? t('nav.expandSidebar') : t('nav.collapseSidebar')}
            onClick={() => setSidebarCollapsed((value) => !value)}
            type="button"
          >
            {sidebarCollapsed ? <PanelLeftOpen size={16} /> : <PanelLeftClose size={16} />}
          </IconButton>
        </div>
        <nav className="workspace-nav">
          {navItems.map((item) => (
            <button
              key={item.value}
              className={view === item.value ? 'active' : ''}
              onClick={() => onViewChange(item.value)}
              type="button"
              title={item.label}
            >
              {item.icon}
              <span className="workspace-nav-label">{item.label}</span>
            </button>
          ))}
        </nav>
        <div className="workspace-sidebar-footer">
          <Button className="workspace-back-button" variant="ghost" onClick={onBack}>
            <ArrowLeft size={16} />
            <span className="workspace-back-label">{t('nav.workspaces')}</span>
          </Button>
        </div>
      </aside>

      <section className={`workspace-main ${topbarHidden ? 'workspace-main-no-topbar' : ''} ${explorerHost ? 'explorer-main-host' : ''}`}>
        {!topbarHidden && (
        <header className="workspace-topbar">
          <div className="row" style={{ minWidth: 0 }}>
            <Badge tone="indigo">{t(viewLabelKey(view))}</Badge>
            <span className="small muted" style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {error || workspace?.paths.root || t('settings.metadata.notLoaded')}
            </span>
          </div>
          <div className="row">
            <HealthBadge health={health} />
            <Button variant="ghost" onClick={() => void refresh()}>
              <RefreshCcw size={15} />
              {t('common.refresh')}
            </Button>
          </div>
        </header>
        )}
        <main className={`workspace-content ${explorerHost ? 'workspace-content-explorer' : ''}`}>
          <Suspense fallback={<div className="workspace-page-loading">{t('common.loading')}</div>}>{page}</Suspense>
        </main>
      </section>
    </div>
  )
}

function viewLabelKey(view: WorkspaceView): MessageKey {
  switch (view) {
    case 'explorer':
      return 'nav.explorer'
    case 'entityTopo':
      return 'nav.entityTopo'
    case 'query':
      return 'nav.query'
    case 'imports':
      return 'nav.imports'
    case 'agent':
      return 'nav.agent'
    case 'settings':
      return 'nav.settings'
    case 'docs':
      return 'nav.apiMap'
    case 'data':
      return 'nav.data'
  }
}
