<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboard, useAgents, useProjects, useTasks } from '@/composables/useApi'
import { useWebSocket } from '@/composables/useWebSocket'

const router = useRouter()

const { data: dashboard, loading: dashLoading, fetch: fetchDashboard } = useDashboard()
const { agents, loading: agentsLoading, fetchAll: fetchAgents } = useAgents()
const { projects, loading: projectsLoading, fetchAll: fetchProjects } = useProjects()
const { tasks, loading: tasksLoading, fetchAll: fetchTasks } = useTasks()

const { on, off, isConnected, connect, lastEvent } = useWebSocket()

const liveEvents = ref([])

const isLoading = computed(() => dashLoading.value || agentsLoading.value || projectsLoading.value || tasksLoading.value)

const activeAgentCount = computed(() => {
  return agents.value.filter(a => {
    const s = (a.status || '').toLowerCase()
    return s === 'working' || s === 'running' || s === 'active'
  }).length
})

const completedTaskCount = computed(() => {
  return tasks.value.filter(t => {
    const s = (t.status || '').toLowerCase()
    return s === 'completed' || s === 'done'
  }).length
})

const systemHealth = computed(() => {
  const d = dashboard.value
  if (d?.systemHealth != null) return d.systemHealth
  if (isConnected.value) return 'Healthy'
  return 'Unknown'
})

const systemHealthPercent = computed(() => {
  const h = systemHealth.value
  if (typeof h === 'number') return h
  const s = String(h).toLowerCase()
  if (s === 'healthy' || s === 'ok' || s === 'good') return 95
  if (s === 'degraded' || s === 'warning') return 60
  if (s === 'critical' || s === 'error' || s === 'down') return 15
  return 0
})

const systemHealthColor = computed(() => {
  const p = systemHealthPercent.value
  if (p >= 80) return 'var(--color-success)'
  if (p >= 40) return 'var(--color-warning)'
  return 'var(--color-error)'
})

const stats = computed(() => [
  {
    label: 'Projects',
    value: projects.value?.length ?? dashboard.value?.projects ?? 0,
    icon: 'projects',
    color: 'var(--color-accent)',
  },
  {
    label: 'Active Agents',
    value: activeAgentCount.value || (dashboard.value?.activeAgents ?? 0),
    icon: 'agents',
    color: 'var(--color-info, #3b82f6)',
  },
  {
    label: 'Tasks Completed',
    value: completedTaskCount.value || (dashboard.value?.tasksDone ?? 0),
    icon: 'tasks',
    color: 'var(--color-success)',
  },
  {
    label: 'System Health',
    value: systemHealth.value,
    icon: 'health',
    color: systemHealthColor.value,
    percent: systemHealthPercent.value,
  },
])

const timeline = computed(() => {
  const apiTimeline = dashboard.value?.timeline ?? []
  return [...liveEvents.value, ...apiTimeline].slice(0, 15)
})

function getEventIcon(type) {
  const t = (type || '').toLowerCase()
  if (t.includes('agent') && t.includes('start')) return 'play'
  if (t.includes('agent') && t.includes('complet')) return 'check'
  if (t.includes('agent') && t.includes('error')) return 'alert'
  if (t.includes('task') && t.includes('creat')) return 'plus'
  if (t.includes('task') && t.includes('complet')) return 'check'
  if (t.includes('deploy')) return 'rocket'
  if (t.includes('build')) return 'build'
  if (t.includes('error') || t.includes('fail')) return 'alert'
  if (t.includes('system')) return 'system'
  return 'event'
}

function getEventColor(type) {
  const t = (type || '').toLowerCase()
  if (t.includes('error') || t.includes('fail')) return 'var(--color-error)'
  if (t.includes('complet') || t.includes('success')) return 'var(--color-success)'
  if (t.includes('start') || t.includes('deploy')) return 'var(--color-accent-bright)'
  if (t.includes('warning')) return 'var(--color-warning)'
  return 'var(--color-text-muted)'
}

function formatTime(ts) {
  if (!ts) return ''
  const d = new Date(ts)
  if (isNaN(d.getTime())) return String(ts)
  const now = new Date()
  const diffMs = now - d
  const diffMin = Math.floor(diffMs / 60000)
  if (diffMin < 1) return 'just now'
  if (diffMin < 60) return `${diffMin}m ago`
  const diffHr = Math.floor(diffMin / 60)
  if (diffHr < 24) return `${diffHr}h ago`
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

// Health ring SVG dash calculation
const healthCircumference = 2 * Math.PI * 18
function healthDashOffset(percent) {
  return healthCircumference - (healthCircumference * Math.max(0, Math.min(100, percent))) / 100
}

function navigateToAgents() { router.push('/agents') }
function navigateToProjects() { router.push('/projects') }

async function handleStartAgent() {
  router.push('/agents')
}

async function handleNewTask() {
  router.push('/projects')
}

async function handleDeploy() {
  router.push('/system')
}

// WebSocket handlers
function onDashboardUpdate(data) {
  if (data.dashboard) {
    dashboard.value = { ...dashboard.value, ...data.dashboard }
  }
}

function onLiveEvent(data) {
  const evt = {
    type: data.type || 'event',
    agent: data.agent?.name || data.agentName || data.agent || '',
    message: data.message || data.description || '',
    timestamp: data.timestamp || new Date().toISOString(),
  }
  liveEvents.value.unshift(evt)
  if (liveEvents.value.length > 20) {
    liveEvents.value = liveEvents.value.slice(0, 20)
  }
}

onMounted(async () => {
  connect()
  on('dashboard_update', onDashboardUpdate)
  on('agent_started', onLiveEvent)
  on('agent_completed', onLiveEvent)
  on('agent_error', onLiveEvent)
  on('task_created', onLiveEvent)
  on('task_completed', onLiveEvent)
  on('deployment', onLiveEvent)
  on('*', (data) => {
    if (data.type && !['dashboard_update'].includes(data.type)) {
      onLiveEvent(data)
    }
  })

  await Promise.allSettled([
    fetchDashboard(),
    fetchAgents(),
    fetchProjects(),
    fetchTasks(),
  ])
})

onUnmounted(() => {
  off('dashboard_update', onDashboardUpdate)
  off('agent_started', onLiveEvent)
  off('agent_completed', onLiveEvent)
  off('agent_error', onLiveEvent)
  off('task_created', onLiveEvent)
  off('task_completed', onLiveEvent)
  off('deployment', onLiveEvent)
})
</script>

<template>
  <div class="dashboard-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">Dashboard</h1>
        <span class="page-subtitle" v-if="isLoading">Loading&hellip;</span>
        <span class="ws-indicator" :class="{ 'ws-indicator--connected': isConnected }" :title="isConnected ? 'Live' : 'Disconnected'" v-else></span>
      </div>
      <div class="header-right">
        <button class="quick-btn quick-btn--accent" @click="handleStartAgent">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><path d="M4 2l10 6-10 6V2z" fill="currentColor"/></svg>
          Start Agent
        </button>
        <button class="quick-btn quick-btn--default" @click="handleNewTask">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><path d="M8 2v12M2 8h12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          New Task
        </button>
        <button class="quick-btn quick-btn--default" @click="handleDeploy">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><path d="M8 2l6 6-6 6M2 8h12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Deploy
        </button>
      </div>
    </header>

    <!-- Stats Row -->
    <section class="stats-grid">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="stat-card"
        :class="{ 'stat-card--clickable': stat.label === 'Projects' || stat.label === 'Active Agents' }"
        @click="stat.label === 'Projects' ? navigateToProjects() : stat.label === 'Active Agents' ? navigateToAgents() : undefined"
      >
        <div class="stat-icon-wrap" :style="{ '--icon-color': stat.color }">
          <!-- Projects icon -->
          <svg v-if="stat.icon === 'projects'" width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M2 4a2 2 0 012-2h4a2 2 0 012 2v2a2 2 0 01-2 2H4a2 2 0 01-2-2V4zm0 8a2 2 0 012-2h4a2 2 0 012 2v4a2 2 0 01-2 2H4a2 2 0 01-2-2v-4zm10-8a2 2 0 012-2h4a2 2 0 012 2v2a2 2 0 01-2 2h-4a2 2 0 01-2-2V4zm0 8a2 2 0 012-2h4a2 2 0 012 2v4a2 2 0 01-2 2h-4a2 2 0 01-2-2v-4z" fill="currentColor"/>
          </svg>
          <!-- Agents icon -->
          <svg v-else-if="stat.icon === 'agents'" width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M10 10a3.5 3.5 0 100-7 3.5 3.5 0 000 7zM3.5 17c0-3.04 2.91-5.5 6.5-5.5s6.5 2.46 6.5 5.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <circle cx="15" cy="6" r="2" fill="currentColor" opacity="0.4"/>
          </svg>
          <!-- Tasks icon -->
          <svg v-else-if="stat.icon === 'tasks'" width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M5.5 10.5l3 3 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <rect x="2" y="2" width="16" height="16" rx="4" stroke="currentColor" stroke-width="1.5"/>
          </svg>
          <!-- Health ring -->
          <svg v-else-if="stat.icon === 'health'" width="40" height="40" viewBox="0 0 42 42">
            <circle cx="21" cy="21" r="18" fill="none" stroke="rgba(255,255,255,0.06)" stroke-width="3"/>
            <circle
              cx="21" cy="21" r="18" fill="none"
              :stroke="systemHealthColor"
              stroke-width="3"
              stroke-linecap="round"
              :stroke-dasharray="healthCircumference"
              :stroke-dashoffset="healthDashOffset(stat.percent ?? 0)"
              transform="rotate(-90 21 21)"
              class="health-ring"
            />
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
      </div>
    </section>

    <!-- Event Timeline -->
    <section class="timeline-section">
      <div class="timeline-header">
        <h2 class="section-title">Event Timeline</h2>
        <span class="event-count" v-if="timeline.length">{{ timeline.length }} recent</span>
      </div>
      <div class="timeline-panel" v-if="timeline.length">
        <div class="timeline-item" v-for="(item, i) in timeline" :key="i">
          <div class="timeline-icon" :style="{ '--event-color': getEventColor(item.type) }">
            <!-- Play -->
            <svg v-if="getEventIcon(item.type) === 'play'" width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M4 2l10 6-10 6V2z" fill="currentColor"/></svg>
            <!-- Check -->
            <svg v-else-if="getEventIcon(item.type) === 'check'" width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M3 8.5l3.5 3.5 6.5-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            <!-- Alert -->
            <svg v-else-if="getEventIcon(item.type) === 'alert'" width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M8 4v5M8 11v1" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><circle cx="8" cy="8" r="6.5" stroke="currentColor" stroke-width="1.5"/></svg>
            <!-- Plus -->
            <svg v-else-if="getEventIcon(item.type) === 'plus'" width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M8 2v12M2 8h12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            <!-- Rocket -->
            <svg v-else-if="getEventIcon(item.type) === 'rocket'" width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M9.5 1.5s2.5 2 2.5 5c0 3-4 6-4 6l-2-2-2-2s3-4 6-4c-3 0-5 2.5-5 2.5M2 14l3-3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            <!-- Build -->
            <svg v-else-if="getEventIcon(item.type) === 'build'" width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M6 10l4-4M5.5 5.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zm8 8a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            <!-- System -->
            <svg v-else-if="getEventIcon(item.type) === 'system'" width="12" height="12" viewBox="0 0 16 16" fill="none"><circle cx="8" cy="8" r="2.5" stroke="currentColor" stroke-width="1.5"/><path d="M8 1.5v2M8 12.5v2M1.5 8h2M12.5 8h2M3.1 3.1l1.4 1.4M11.5 11.5l1.4 1.4M3.1 12.9l1.4-1.4M11.5 4.5l1.4-1.4" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/></svg>
            <!-- Default event -->
            <svg v-else width="12" height="12" viewBox="0 0 16 16" fill="none"><circle cx="8" cy="8" r="3" fill="currentColor"/></svg>
          </div>
          <div class="timeline-body">
            <div class="timeline-top-row">
              <span class="timeline-agent" v-if="item.agent">{{ item.agent }}</span>
              <span class="timeline-desc">{{ item.message || item.event || item.description || item.type }}</span>
            </div>
            <span class="timeline-time">{{ formatTime(item.timestamp || item.time) }}</span>
          </div>
        </div>
      </div>
      <div class="empty-state" v-else>
        <svg width="32" height="32" viewBox="0 0 32 32" fill="none" class="empty-icon">
          <circle cx="16" cy="16" r="12" stroke="currentColor" stroke-width="1.5" opacity="0.3"/>
          <path d="M16 10v6l4 2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" opacity="0.4"/>
        </svg>
        <span class="empty-text">No recent events</span>
        <span class="empty-sub">Events will appear here as agents work</span>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* ── Page ── */
.dashboard-page {
  padding: var(--space-8, 32px);
  max-width: 1160px;
  margin: 0 auto;
  min-height: 100vh;
  font-family: var(--font-primary, 'Inter', system-ui, sans-serif);
  color: var(--color-text-primary, #f7f8f8);
}

/* ── Header ── */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-8, 32px);
  gap: var(--space-4, 16px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3, 12px);
}

.page-title {
  font-family: var(--font-primary, 'Inter', system-ui, sans-serif);
  font-size: var(--font-h2, 24px);
  font-weight: 600;
  color: var(--color-text-primary, #f7f8f8);
  margin: 0;
  letter-spacing: -0.3px;
}

.page-subtitle {
  font-size: var(--font-caption, 13px);
  color: var(--color-text-muted, #8a8f98);
}

.ws-indicator {
  width: 7px;
  height: 7px;
  border-radius: var(--radius-full, 9999px);
  background: var(--color-error, #ef4444);
  flex-shrink: 0;
  transition: background var(--transition-normal, 250ms ease);
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.ws-indicator--connected {
  background: var(--color-success, #10b981);
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
  animation: ws-pulse 3s ease-in-out infinite;
}

@keyframes ws-pulse {
  0%, 100% { box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2); }
  50% { box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.1); }
}

/* ── Quick Action Buttons ── */
.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2, 8px);
}

.quick-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-primary, 'Inter', system-ui, sans-serif);
  font-size: var(--font-caption, 13px);
  font-weight: 500;
  padding: 7px 14px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--border-standard, rgba(255, 255, 255, 0.08));
  cursor: pointer;
  transition: all var(--transition-fast, 150ms ease);
  white-space: nowrap;
  line-height: 1;
}

.quick-btn svg {
  flex-shrink: 0;
}

.quick-btn--accent {
  background: var(--color-accent, #5e6ad2);
  color: #fff;
  border-color: transparent;
}

.quick-btn--accent:hover {
  background: var(--color-accent-hover, #828fff);
  border-color: transparent;
  transform: translateY(-0.5px);
  box-shadow: 0 2px 8px rgba(94, 106, 210, 0.3);
}

.quick-btn--default {
  background: rgba(255, 255, 255, 0.02);
  color: var(--color-text-secondary, #d0d6e0);
}

.quick-btn--default:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: var(--border-strong, rgba(255, 255, 255, 0.12));
  transform: translateY(-0.5px);
}

.quick-btn:active {
  transform: translateY(0) scale(0.98);
}

/* ── Stats Grid ── */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4, 16px);
  margin-bottom: var(--space-8, 32px);
}

.stat-card {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-lg, 8px);
  padding: var(--space-5, 20px);
  display: flex;
  align-items: center;
  gap: var(--space-4, 16px);
  transition: all var(--transition-fast, 150ms ease);
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.06), transparent);
  opacity: 0;
  transition: opacity var(--transition-fast, 150ms ease);
}

.stat-card:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.12);
  transform: translateY(-1px);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-card--clickable {
  cursor: pointer;
}

.stat-card--clickable:hover {
  border-color: rgba(255, 255, 255, 0.15);
}

/* ── Stat Icon ── */
.stat-icon-wrap {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-xl, 12px);
  background: rgba(255, 255, 255, 0.04);
  color: var(--icon-color, var(--color-accent, #5e6ad2));
  flex-shrink: 0;
}

/* ── Health Ring ── */
.health-ring {
  transition: stroke-dashoffset 1s ease;
}

/* ── Stat Info ── */
.stat-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-value {
  font-family: var(--font-primary, 'Inter', system-ui, sans-serif);
  font-size: var(--font-h3, 20px);
  font-weight: 600;
  color: var(--color-text-primary, #f7f8f8);
  line-height: 1.2;
  letter-spacing: -0.3px;
}

.stat-label {
  font-size: var(--font-label, 12px);
  font-weight: 500;
  color: var(--color-text-muted, #8a8f98);
  margin-top: 2px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* ── Timeline Section ── */
.timeline-section {
  background: rgba(255, 255, 255, 0.015);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-lg, 8px);
  padding: var(--space-5, 20px);
}

.timeline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4, 16px);
}

.section-title {
  font-family: var(--font-primary, 'Inter', system-ui, sans-serif);
  font-size: var(--font-small, 15px);
  font-weight: 600;
  color: var(--color-text-secondary, #d0d6e0);
  margin: 0;
  letter-spacing: -0.1px;
}

.event-count {
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: var(--font-micro, 11px);
  color: var(--color-text-muted, #8a8f98);
  background: rgba(255, 255, 255, 0.04);
  padding: 3px 8px;
  border-radius: var(--radius-full, 9999px);
}

/* ── Timeline Panel ── */
.timeline-panel {
  display: flex;
  flex-direction: column;
}

.timeline-item {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3, 12px);
  padding: var(--space-3, 12px) 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  transition: background var(--transition-fast, 150ms ease);
  border-radius: var(--radius-md, 6px);
  padding-left: var(--space-2, 8px);
  padding-right: var(--space-2, 8px);
}

.timeline-item:last-child {
  border-bottom: none;
}

.timeline-item:hover {
  background: rgba(255, 255, 255, 0.02);
}

/* ── Timeline Icon ── */
.timeline-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md, 6px);
  background: rgba(255, 255, 255, 0.04);
  color: var(--event-color, var(--color-text-muted, #8a8f98));
  flex-shrink: 0;
  margin-top: 1px;
}

/* ── Timeline Body ── */
.timeline-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.timeline-top-row {
  display: flex;
  align-items: baseline;
  gap: var(--space-2, 8px);
  flex-wrap: wrap;
}

.timeline-agent {
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: var(--font-caption, 13px);
  font-weight: 600;
  color: var(--color-accent-bright, #7170ff);
  white-space: nowrap;
}

.timeline-desc {
  font-size: var(--font-caption, 13px);
  color: var(--color-text-secondary, #d0d6e0);
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.timeline-time {
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: var(--font-micro, 11px);
  color: var(--color-text-muted, #8a8f98);
}

/* ── Empty State ── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12, 48px) 0;
  gap: var(--space-2, 8px);
}

.empty-icon {
  color: var(--color-text-muted, #8a8f98);
  margin-bottom: var(--space-2, 8px);
}

.empty-text {
  font-size: var(--font-small, 15px);
  font-weight: 500;
  color: var(--color-text-muted, #8a8f98);
}

.empty-sub {
  font-size: var(--font-caption, 13px);
  color: var(--color-text-subtle, #62666d);
}

/* ── Responsive ── */
@media (max-width: 768px) {
  .dashboard-page {
    padding: var(--space-4, 16px);
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-4, 16px);
  }

  .header-right {
    width: 100%;
    overflow-x: auto;
    padding-bottom: var(--space-1, 4px);
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-card {
    padding: var(--space-4, 16px);
  }

  .stat-value {
    font-size: var(--font-body, 16px);
  }
}

@media (max-width: 480px) {
  .quick-btn span:last-child {
    display: none;
  }

  .quick-btn {
    padding: 8px 10px;
  }

  .timeline-top-row {
    flex-direction: column;
    gap: 2px;
  }
}
</style>