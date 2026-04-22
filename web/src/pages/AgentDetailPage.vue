<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAgents } from '@/composables/useApi'
import { useWebSocket } from '@/composables/useWebSocket'

const route = useRoute()
const agentId = computed(() => route.params.id)

const { fetchOne, startAgent } = useAgents()
const { on, off, isConnected, connect } = useWebSocket()

const agent = ref(null)
const terminalLines = ref([])
const loading = ref(true)
const terminalRef = ref(null)
const stopping = ref(false)
const restarting = ref(false)

function normalizeStatus(status) {
  const s = (status || '').toLowerCase()
  if (s === 'working' || s === 'running') return 'working'
  if (s === 'error' || s === 'failed') return 'error'
  if (s === 'completed' || s === 'done') return 'completed'
  return 'idle'
}

function badgeClass(status) {
  return `badge--${normalizeStatus(status)}`
}

async function loadAgent() {
  try {
    agent.value = await fetchOne(agentId.value)
  } catch {
    agent.value = null
  } finally {
    loading.value = false
  }
}

function scrollToBottom() {
  nextTick(() => {
    const el = terminalRef.value
    if (el) el.scrollTop = el.scrollHeight
  })
}

// WebSocket handlers
function onAgentOutput(data) {
  if (data.agentId === agentId.value || data.id === agentId.value) {
    const text = data.message || data.text || data.output || ''
    if (text) {
      const lines = text.split('\n').filter(l => l.length > 0)
      terminalLines.value.push(...lines)
      if (terminalLines.value.length > 500) {
        terminalLines.value = terminalLines.value.slice(-500)
      }
      scrollToBottom()
    }
  }
}

function onAgentStarted(data) {
  const a = data.agent || data
  if (a.id === agentId.value || a.agentId === agentId.value) {
    agent.value = { ...agent.value, ...a, status: 'working' }
  }
}

function onAgentCompleted(data) {
  const a = data.agent || data
  if (a.id === agentId.value || a.agentId === agentId.value) {
    agent.value = { ...agent.value, ...a, status: 'completed' }
  }
}

function onAgentError(data) {
  const a = data.agent || data
  if (a.id === agentId.value || a.agentId === agentId.value) {
    agent.value = { ...agent.value, ...a, status: 'error' }
  }
}

async function handleStop() {
  stopping.value = true
  try {
    await fetchOne(agentId.value)
    agent.value = await fetchOne(agentId.value)
  } catch { /* ignore */ }
  stopping.value = false
}

async function handleRestart() {
  restarting.value = true
  terminalLines.value = []
  try {
    const name = agent.value?.name
    const task = agent.value?.task || agent.value?.currentTask
    if (name && task) {
      await startAgent(name, task)
    }
    agent.value = await fetchOne(agentId.value)
  } catch { /* ignore */ }
  restarting.value = false
}

// Re-fetch agent if route id changes
watch(agentId, (newId) => {
  if (newId) {
    loading.value = true
    agent.value = null
    terminalLines.value = []
    loadAgent()
  }
})

onMounted(async () => {
  await loadAgent()
  connect()
  on('agent_output', onAgentOutput)
  on('agent_started', onAgentStarted)
  on('agent_completed', onAgentCompleted)
  on('agent_error', onAgentError)
})

onUnmounted(() => {
  off('agent_output', onAgentOutput)
  off('agent_started', onAgentStarted)
  off('agent_completed', onAgentCompleted)
  off('agent_error', onAgentError)
})
</script>

<template>
  <div class="agent-detail-page">
    <div class="loading-state" v-if="loading">
      <span class="loading-text">Loading agent…</span>
    </div>

    <div class="not-found" v-else-if="!agent">
      <span class="not-found-text">Agent not found</span>
    </div>

    <template v-else>
      <header class="page-header">
        <div class="header-left">
          <span class="agent-emoji">{{ agent.emoji || '🤖' }}</span>
          <div class="header-text">
            <div class="title-row">
              <h1 class="page-title">{{ agent.name }}</h1>
              <span class="status-badge" :class="badgeClass(agent.status)">
                <span class="status-dot"></span>
                {{ agent.status || 'idle' }}
              </span>
            </div>
            <span class="agent-model">{{ agent.model || 'Unknown model' }}</span>
          </div>
        </div>
        <span class="ws-indicator" :class="{ 'ws-indicator--connected': isConnected }"></span>
      </header>

      <section class="info-section">
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">Model</span>
            <span class="info-value">{{ agent.model || '—' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Channel</span>
            <span class="info-value">{{ agent.channel || '—' }}</span>
          </div>
          <div class="info-item" v-if="agent.task || agent.currentTask">
            <span class="info-label">Current Task</span>
            <span class="info-value info-value--task">{{ agent.task || agent.currentTask }}</span>
          </div>
          <div class="info-item" v-if="agent.uptime">
            <span class="info-label">Uptime</span>
            <span class="info-value">{{ agent.uptime }}</span>
          </div>
          <div class="info-item" v-if="agent.tokensUsed ?? agent.tokens">
            <span class="info-label">Tokens Used</span>
            <span class="info-value">{{ agent.tokensUsed ?? agent.tokens }}</span>
          </div>
        </div>
      </section>

      <section class="terminal-section">
        <div class="terminal-header">
          <h2 class="section-title">Live Output</h2>
          <span class="line-count" v-if="terminalLines.length">{{ terminalLines.length }} lines</span>
        </div>
        <div class="terminal">
          <div class="terminal-scroll" ref="terminalRef">
            <div class="terminal-line" v-for="(line, i) in terminalLines" :key="i">{{ line }}</div>
            <div class="terminal-cursor" v-if="normalizeStatus(agent.status) === 'working'">
              <span class="cursor-block"></span>
            </div>
            <div class="terminal-empty" v-if="!terminalLines.length">
              Waiting for output…
            </div>
          </div>
        </div>
      </section>

      <section class="actions-section">
        <button
          class="action-btn action-btn--stop"
          @click="handleStop"
          :disabled="stopping || normalizeStatus(agent.status) === 'idle'"
        >
          <span class="action-btn-icon">■</span>
          {{ stopping ? 'Stopping…' : 'Stop' }}
        </button>
        <button
          class="action-btn action-btn--restart"
          @click="handleRestart"
          :disabled="restarting"
        >
          <span class="action-btn-icon">⟳</span>
          {{ restarting ? 'Restarting…' : 'Restart' }}
        </button>
      </section>
    </template>
  </div>
</template>

<style scoped>
.agent-detail-page {
  padding: var(--space-8);
  max-width: 1120px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-6);
  padding-bottom: var(--space-5);
  border-bottom: 1px solid var(--border-standard);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.agent-emoji {
  font-size: 40px;
  line-height: 1;
}

.header-text {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.title-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.page-title {
  font-family: var(--font-primary);
  font-size: var(--font-h2);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.ws-indicator {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--color-error);
  flex-shrink: 0;
  transition: background var(--transition-normal);
}

.ws-indicator--connected {
  background: var(--color-success);
}

/* Status badge */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  font-weight: 500;
  border-radius: var(--radius-full);
  padding: 2px var(--space-2) 2px 6px;
  text-transform: capitalize;
  white-space: nowrap;
  flex-shrink: 0;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: var(--radius-full);
  flex-shrink: 0;
}

.badge--idle {
  color: var(--color-text-muted);
  background: rgba(138, 143, 152, 0.12);
}
.badge--idle .status-dot {
  background: var(--color-text-muted);
}

.badge--working {
  color: var(--color-info);
  background: rgba(59, 130, 246, 0.12);
}
.badge--working .status-dot {
  background: var(--color-info);
  animation: pulse-dot 2s ease-in-out infinite;
}

.badge--completed {
  color: var(--color-success);
  background: rgba(16, 185, 129, 0.12);
}
.badge--completed .status-dot {
  background: var(--color-success);
}

.badge--error {
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.12);
}
.badge--error .status-dot {
  background: var(--color-error);
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.85); }
}

.agent-model {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

/* Info section */
.info-section {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  margin-bottom: var(--space-6);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: var(--space-5);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.info-label {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  font-weight: 500;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
  word-break: break-word;
}

.info-value--task {
  font-size: var(--font-small);
  color: var(--color-text-primary);
  line-height: 1.5;
}

/* Terminal section */
.terminal-section {
  margin-bottom: var(--space-6);
}

.terminal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.section-title {
  font-family: var(--font-primary);
  font-size: var(--font-small);
  font-weight: 600;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0;
}

.line-count {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  color: var(--color-text-muted);
}

.terminal {
  background: #0a0a0a;
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  overflow: hidden;
}

.terminal-scroll {
  font-family: var(--font-mono);
  font-size: var(--font-caption);
  color: #4ade80;
  line-height: 1.7;
  max-height: 400px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

/* Custom scrollbar for terminal */
.terminal-scroll::-webkit-scrollbar {
  width: 6px;
}

.terminal-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.terminal-scroll::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-full);
}

.terminal-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}

.terminal-line {
  min-height: 1.7em;
}

.terminal-cursor {
  display: inline-flex;
  align-items: center;
  height: 1.7em;
}

.cursor-block {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #4ade80;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.terminal-empty {
  color: var(--color-text-subtle);
  font-style: italic;
}

/* Actions section */
.actions-section {
  display: flex;
  gap: var(--space-3);
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  font-weight: 500;
  padding: var(--space-3) var(--space-6);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-standard);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.action-btn-icon {
  font-size: 14px;
  line-height: 1;
}

.action-btn--stop {
  background: rgba(239, 68, 68, 0.12);
  color: var(--color-error);
  border-color: rgba(239, 68, 68, 0.2);
}

.action-btn--stop:not(:disabled):hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.35);
}

.action-btn--restart {
  background: var(--color-surface);
  color: var(--color-text-secondary);
}

.action-btn--restart:not(:disabled):hover {
  background: var(--color-hover);
  border-color: var(--border-strong);
}

/* Loading / Not found */
.loading-state,
.not-found {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-16) 0;
}

.loading-text,
.not-found-text {
  font-family: var(--font-primary);
  font-size: var(--font-small);
  color: var(--color-text-muted);
}

/* Mobile */
@media (max-width: 768px) {
  .agent-detail-page {
    padding: var(--space-4);
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }

  .header-left {
    width: 100%;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .terminal-scroll {
    max-height: 280px;
  }

  .actions-section {
    flex-direction: column;
  }

  .action-btn {
    justify-content: center;
  }
}
</style>