<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAgents } from '@/composables/useApi'
import { useWebSocket } from '@/composables/useWebSocket'

const router = useRouter()
const { agents, loading, fetchAll } = useAgents()
const { on, off, isConnected, connect } = useWebSocket()

const activeFilter = ref('all')
const hoveredAgent = ref(null)

const filters = [
  { key: 'all', label: 'All', icon: '⚡' },
  { key: 'working', label: 'Active', icon: '🔥' },
  { key: 'standby', label: 'Standby', icon: '💤' },
  { key: 'error', label: 'Error', icon: '⚠️' },
]

const filteredAgents = computed(() => {
  if (activeFilter.value === 'all') return agents.value
  return agents.value.filter(a => norm(a.status) === activeFilter.value)
})

const workingCount = computed(() => agents.value.filter(a => norm(a.status) === 'working').length)
const standbyCount = computed(() => agents.value.filter(a => norm(a.status) === 'standby').length)

function norm(s) {
  s = (s || '').toLowerCase()
  if (s === 'working' || s === 'running') return 'working'
  if (s === 'error' || s === 'failed') return 'error'
  if (s === 'completed' || s === 'done') return 'completed'
  if (s === 'standby') return 'standby'
  return 'idle'
}

function providerInfo(channel) {
  const c = (channel || '').toLowerCase()
  if (c === 'ollama') return { icon: '🦙', label: 'Ollama', color: 'purple' }
  if (c === 'cloud') return { icon: '☁️', label: 'Cloud', color: 'blue' }
  if (c === 'telegram') return { icon: '📱', label: 'Telegram', color: 'sky' }
  if (c === 'discord') return { icon: '💬', label: 'Discord', color: 'indigo' }
  return { icon: '🖥️', label: 'Local', color: 'gray' }
}

function shortModel(m) {
  if (!m) return '—'
  return m.replace('ollama/', '').replace(':latest', '')
}

function formatTime(t) {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return ''
  const h = d.getHours(), m = d.getMinutes()
  const ampm = h >= 12 ? 'PM' : 'AM'
  return `${h % 12 || 12}:${String(m).padStart(2, '0')} ${ampm}`
}

function fmtDuration(ms) {
  if (!ms) return ''
  const s = Math.floor(ms / 1000)
  if (s < 60) return `${s}s`
  const m = Math.floor(s / 60)
  if (m < 60) return `${m}m ${s % 60}s`
  return `${Math.floor(m / 60)}h ${m % 60}m`
}

function fmtTokens(n) {
  if (!n) return ''
  if (n >= 1000000) return `${(n / 1000000).toFixed(1)}M`
  if (n >= 1000) return `${(n / 1000).toFixed(1)}K`
  return `${n}`
}

function go(id) { router.push(`/agents/${id}`) }

function onAgentsUpdate(data) {
  if (Array.isArray(data.agents)) agents.value = data.agents
  else if (data.agent) {
    const i = agents.value.findIndex(a => a.id === data.agent.id)
    if (i !== -1) agents.value[i] = { ...agents.value[i], ...data.agent }
    else agents.value.push(data.agent)
  }
}

function onAgentEvent(data) {
  const agent = data.agent || data
  const i = agents.value.findIndex(a => a.id === agent.id)
  if (i !== -1) agents.value[i] = { ...agents.value[i], ...agent }
  else agents.value.push({ ...agent })
}

onMounted(() => {
  fetchAll()
  connect()
  on('agents_update', onAgentsUpdate)
  on('agent_started', onAgentEvent)
  on('agent_completed', onAgentEvent)
  on('agent_error', onAgentEvent)
  on('agent_status_update', onAgentEvent)
})

onUnmounted(() => {
  off('agents_update', onAgentsUpdate)
  off('agent_started', onAgentEvent)
  off('agent_completed', onAgentEvent)
  off('agent_error', onAgentEvent)
  off('agent_status_update', onAgentEvent)
})
</script>

<template>
  <div class="command-center">
    <!-- Hero Header -->
    <div class="hero">
      <div class="hero-glow"></div>
      <div class="hero-content">
        <div class="hero-left">
          <h1 class="hero-title">Command Center</h1>
          <p class="hero-subtitle">{{ agents.length }} agents in the pack</p>
        </div>
        <div class="hero-stats">
          <div class="stat-chip stat-active" v-if="workingCount">
            <span class="stat-dot"></span>
            {{ workingCount }} active
          </div>
          <div class="stat-chip stat-standby" v-if="standbyCount">
            <span class="stat-dot"></span>
            {{ standbyCount }} standby
          </div>
          <div class="stat-chip stat-live" :class="{ offline: !isConnected }">
            <span class="stat-pulse"></span>
            {{ isConnected ? 'LIVE' : 'OFFLINE' }}
          </div>
        </div>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="controls">
      <div class="filter-group">
        <button
          v-for="f in filters" :key="f.key"
          class="filter-pill" :class="{ active: activeFilter === f.key }"
          @click="activeFilter = f.key"
        >
          <span class="filter-em">{{ f.icon }}</span>
          {{ f.label }}
        </button>
      </div>
    </div>

    <!-- Agent List — Table-like command center -->
    <div class="agent-list" v-if="!loading && filteredAgents.length">
      <div
        v-for="agent in filteredAgents"
        :key="agent.id"
        class="agent-row"
        :class="[`agent-row--${norm(agent.status)}`, { hovered: hoveredAgent === agent.id }]"
        @click="go(agent.id)"
        @mouseenter="hoveredAgent = agent.id"
        @mouseleave="hoveredAgent = null"
      >
        <!-- Accent bar -->
        <div class="row-accent" :class="`accent--${norm(agent.status)}`"></div>

        <!-- Avatar -->
        <div class="row-avatar" :class="`avatar--${norm(agent.status)}`">
          <span class="avatar-emoji">{{ agent.emoji || '🤖' }}</span>
          <span class="avatar-status-dot" :class="norm(agent.status)"></span>
        </div>

        <!-- Name + Task -->
        <div class="row-main">
          <div class="row-name-row">
            <span class="row-name">{{ agent.name }}</span>
            <span class="row-id" v-if="hoveredAgent === agent.id">#{{ agent.id }}</span>
          </div>
          <span class="row-task" v-if="agent.task || agent.currentTask">
            {{ agent.task || agent.currentTask }}
          </span>
          <span class="row-task row-task--none" v-else>Awaiting orders</span>
        </div>

        <!-- Model + Provider -->
        <div class="row-model">
          <span class="provider-tag" :class="`tag--${providerInfo(agent.channel).color}`">
            {{ providerInfo(agent.channel).icon }} {{ shortModel(agent.model) }}
          </span>
        </div>

        <!-- Metrics -->
        <div class="row-metrics">
          <span class="metric-chip" v-if="agent.tokensIn || agent.tokensOut" :title="`In: ${agent.tokensIn} / Out: ${agent.tokensOut}`">
            📊 {{ fmtTokens((agent.tokensIn || 0) + (agent.tokensOut || 0)) }} tok
          </span>
          <span class="metric-chip" v-if="agent.durationMs" :title="`${agent.durationMs}ms`">
            ⏱ {{ fmtDuration(agent.durationMs) }}
          </span>
        </div>

        <!-- Status -->
        <div class="row-status">
          <span class="status-tag" :class="`tag--${norm(agent.status)}`">
            <span class="tag-dot" :class="norm(agent.status)"></span>
            {{ norm(agent.status) }}
          </span>
        </div>

        <!-- Time -->
        <div class="row-time">
          {{ formatTime(agent.updatedAt) }}
        </div>

        <!-- Arrow -->
        <div class="row-arrow">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M6 3l5 5-5 5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div class="empty" v-else-if="!loading">
      <div class="empty-icon">🤖</div>
      <div class="empty-text" v-if="agents.length">No agents match "{{ activeFilter }}"</div>
      <div class="empty-text" v-else>No agents connected</div>
    </div>

    <!-- Loading -->
    <div class="loading" v-else>
      <div class="loader"></div>
      <span>Scanning agents…</span>
    </div>
  </div>
</template>

<style scoped>
.command-center {
  max-width: 1100px;
  margin: 0 auto;
}

/* ── Hero ── */
.hero {
  position: relative;
  padding: 40px 32px 28px;
  border-bottom: 1px solid var(--border-subtle);
  overflow: hidden;
}

.hero-glow {
  position: absolute;
  top: -80px;
  left: 50%;
  transform: translateX(-50%);
  width: 600px;
  height: 200px;
  background: radial-gradient(ellipse, rgba(94, 106, 210, 0.12) 0%, transparent 70%);
  pointer-events: none;
  animation: hero-breathe 6s ease-in-out infinite;
}

@keyframes hero-breathe {
  0%, 100% { opacity: 0.6; transform: translateX(-50%) scale(1); }
  50% { opacity: 1; transform: translateX(-50%) scale(1.1); }
}

.hero-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 1;
}

.hero-title {
  font-family: var(--font-primary);
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.5px;
}

.hero-subtitle {
  font-size: var(--font-caption);
  color: var(--color-text-muted);
  margin: 4px 0 0;
}

.hero-stats {
  display: flex;
  gap: 8px;
}

.stat-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 5px 12px;
  border-radius: var(--radius-full);
  border: 1px solid transparent;
}

.stat-active {
  color: #7170ff;
  background: rgba(113, 112, 255, 0.1);
  border-color: rgba(113, 112, 255, 0.2);
}

.stat-standby {
  color: var(--color-text-muted);
  background: rgba(138, 143, 152, 0.08);
  border-color: rgba(138, 143, 152, 0.12);
}

.stat-live {
  color: var(--color-success);
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
}

.stat-live.offline {
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
}

.stat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.stat-pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  animation: stat-blink 2s ease-in-out infinite;
}

@keyframes stat-blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

/* ── Controls ── */
.controls {
  padding: 16px 32px;
  border-bottom: 1px solid var(--border-subtle);
}

.filter-group {
  display: flex;
  gap: 6px;
}

.filter-pill {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-primary);
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-muted);
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--radius-full);
  padding: 4px 14px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.filter-pill:hover {
  color: var(--color-text-secondary);
  background: var(--color-surface);
}

.filter-pill.active {
  color: var(--color-text-primary);
  background: var(--color-elevated);
  border-color: var(--border-strong);
}

.filter-em {
  font-size: 11px;
}

/* ── Agent List ── */
.agent-list {
  padding: 8px 0;
}

.agent-row {
  display: grid;
  grid-template-columns: 3px 48px 1fr 160px 140px 100px 80px 24px;
  align-items: center;
  gap: 0;
  padding: 14px 28px 14px 0;
  margin: 0 16px;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
}

.agent-row:hover {
  background: var(--color-surface);
}

.agent-row:hover .row-arrow {
  opacity: 1;
  color: var(--color-text-muted);
}

.agent-row:hover .row-id {
  opacity: 1;
}

/* Accent bar */
.row-accent {
  width: 3px;
  height: 40px;
  border-radius: 2px;
  background: transparent;
  transition: all 0.2s ease;
  margin-right: 12px;
  flex-shrink: 0;
}

.accent--working {
  background: linear-gradient(180deg, #7170ff, #5e6ad2);
  box-shadow: 0 0 8px rgba(113, 112, 255, 0.4);
  animation: accent-glow 2s ease-in-out infinite;
}

.accent--standby {
  background: var(--color-text-subtle);
  opacity: 0.4;
}

.accent--error {
  background: var(--color-error);
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.3);
}

.accent--completed {
  background: var(--color-success);
}

.accent--idle {
  background: transparent;
}

@keyframes accent-glow {
  0%, 100% { box-shadow: 0 0 4px rgba(113, 112, 255, 0.3); }
  50% { box-shadow: 0 0 12px rgba(113, 112, 255, 0.6); }
}

/* Avatar */
.row-avatar {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-elevated);
  border-radius: 12px;
  font-size: 20px;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.avatar--working {
  background: linear-gradient(135deg, rgba(94, 106, 210, 0.2), rgba(113, 112, 255, 0.1));
  box-shadow: 0 0 16px rgba(94, 106, 210, 0.15);
}

.avatar-emoji {
  line-height: 1;
}

.avatar-status-dot {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid var(--color-canvas);
  z-index: 1;
}

.avatar-status-dot.working {
  background: #7170ff;
  animation: dot-pulse 1.5s ease-in-out infinite;
}

.avatar-status-dot.standby {
  background: var(--color-text-subtle);
}

.avatar-status-dot.error {
  background: var(--color-error);
}

.avatar-status-dot.completed {
  background: var(--color-success);
}

.avatar-status-dot.idle {
  background: var(--color-text-subtle);
  opacity: 0.5;
}

@keyframes dot-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.8); }
}

/* Main */
.row-main {
  min-width: 0;
  padding-left: 12px;
}

.row-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.row-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
}

.row-id {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--color-text-subtle);
  opacity: 0;
  transition: opacity 0.15s ease;
}

.row-task {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
  max-width: 300px;
}

.row-task--none {
  color: var(--color-text-subtle);
  font-style: italic;
}

/* Model */
.row-model {
  display: flex;
  align-items: center;
}

.provider-tag {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: var(--radius-md);
  white-space: nowrap;
}

.tag--purple {
  color: #c4b5fd;
  background: rgba(167, 139, 250, 0.12);
  border: 1px solid rgba(167, 139, 250, 0.15);
}

.tag--blue {
  color: #93c5fd;
  background: rgba(96, 165, 250, 0.12);
  border: 1px solid rgba(96, 165, 250, 0.15);
}

.tag--sky {
  color: #7dd3fc;
  background: rgba(56, 189, 248, 0.12);
  border: 1px solid rgba(56, 189, 248, 0.15);
}

.tag--indigo {
  color: #a5b4fc;
  background: rgba(165, 180, 252, 0.12);
  border: 1px solid rgba(165, 180, 252, 0.15);
}

.tag--gray {
  color: var(--color-text-muted);
  background: rgba(138, 143, 152, 0.08);
  border: 1px solid rgba(138, 143, 152, 0.1);
}

/* Metrics */
.row-metrics {
  display: flex;
  align-items: center;
  gap: 6px;
}

.metric-chip {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 500;
  color: var(--color-text-muted);
  background: var(--color-canvas);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  white-space: nowrap;
}

/* Status */
.row-status {
  display: flex;
  align-items: center;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  padding: 3px 10px;
  border-radius: var(--radius-full);
}

.tag-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  flex-shrink: 0;
}

.tag--working {
  color: #7170ff;
  background: rgba(113, 112, 255, 0.1);
}
.tag--working .tag-dot { background: #7170ff; animation: dot-pulse 1.5s ease-in-out infinite; }

.tag--standby {
  color: var(--color-text-muted);
  background: rgba(138, 143, 152, 0.08);
}
.tag--standby .tag-dot { background: var(--color-text-subtle); }

.tag--idle {
  color: var(--color-text-subtle);
  background: rgba(138, 143, 152, 0.05);
}
.tag--idle .tag-dot { background: var(--color-text-subtle); opacity: 0.5; }

.tag--error {
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.1);
}
.tag--error .tag-dot { background: var(--color-error); }

.tag--completed {
  color: var(--color-success);
  background: rgba(16, 185, 129, 0.1);
}
.tag--completed .tag-dot { background: var(--color-success); }

/* Time */
.row-time {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--color-text-subtle);
  text-align: right;
  tabular-nums: true;
  font-variant-numeric: tabular-nums;
}

/* Arrow */
.row-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  color: var(--color-text-subtle);
  transition: all 0.15s ease;
}

/* ── Empty & Loading ── */
.empty, .loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 80px 0;
}

.empty-icon {
  font-size: 40px;
}

.empty-text, .loading span {
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.loader {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-standard);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ── Responsive ── */
@media (max-width: 900px) {
  .hero { padding: 24px 16px 20px; }
  .hero-content { flex-direction: column; align-items: flex-start; gap: 12px; }
  .hero-stats { width: 100%; flex-wrap: wrap; }
  .controls { padding: 12px 16px; }

  .agent-row {
    grid-template-columns: 3px 40px 1fr auto 24px;
    padding: 12px 12px 12px 0;
    margin: 0 8px;
  }
  .row-model, .row-metrics, .row-time { display: none; }
  .row-main { padding-left: 8px; }
}

@media (max-width: 480px) {
  .agent-row {
    grid-template-columns: 3px 36px 1fr 24px;
  }
  .row-status { display: none; }
}
</style>