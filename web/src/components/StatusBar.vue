<template>
  <div class="statusbar">
    <div class="statusbar-inner">
      <div class="status-items" v-if="activeAgents.length > 0">
        <div
          v-for="agent in activeAgents.slice(0, 3)"
          :key="agent.id"
          class="status-item"
        >
          <span class="status-dot" :class="agent.status"></span>
          <span class="status-icon">{{ agent.emoji || '🤖' }}</span>
          <span class="status-name">{{ agent.name }}</span>
          <span class="status-text">{{ statusLabel(agent.status) }}</span>
          <span class="status-time" v-if="agent.task">{{ truncate(agent.task, 30) }}</span>
        </div>
      </div>
      <div class="status-idle" v-else>
        <span class="status-dot standby"></span>
        <span class="status-text">{{ agents.length }} agents on standby</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'

const { connect, on, off } = useWebSocket()
const agents = ref([])
const activeAgents = computed(() => agents.value.filter(a => a.status === 'working'))

function statusLabel(status) {
  const labels = {
    idle: 'idle',
    standby: 'standby',
    working: 'working...',
    completed: 'done',
    error: 'error',
  }
  return labels[status] || status
}

function truncate(str, len) {
  if (!str) return ''
  return str.length > len ? str.slice(0, len) + '...' : str
}

function handleAgentsUpdate(data) {
  if (data.agents) {
    agents.value = data.agents
  }
}

function handleAgentEvent(data) {
  if (data.agent) {
    const idx = agents.value.findIndex(a => a.id === data.agent.id)
    if (idx >= 0) {
      agents.value[idx] = data.agent
    } else {
      agents.value.push(data.agent)
    }
  }
}

onMounted(() => {
  connect()
  on('agents_update', handleAgentsUpdate)
  on('agent_started', handleAgentEvent)
  on('agent_output', handleAgentEvent)
  on('agent_completed', handleAgentEvent)
  on('agent_error', handleAgentEvent)
})

onUnmounted(() => {
  off('agents_update', handleAgentsUpdate)
  off('agent_started', handleAgentEvent)
  off('agent_output', handleAgentEvent)
  off('agent_completed', handleAgentEvent)
  off('agent_error', handleAgentEvent)
})
</script>

<style scoped>
.statusbar {
  height: var(--statusbar-height);
  background: var(--color-panel);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  padding: 0 var(--space-4);
  flex-shrink: 0;
  overflow: hidden;
}

.statusbar-inner {
  width: 100%;
  overflow: hidden;
}

.status-items {
  display: flex;
  gap: var(--space-4);
  overflow-x: auto;
}

.status-item {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--font-micro);
  color: var(--color-text-muted);
  white-space: nowrap;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-dot.idle { background: var(--color-text-subtle); }
.status-dot.standby { background: var(--color-warning); animation: pulse 2s infinite; }
.status-dot.working { background: var(--color-accent-bright); animation: pulse 2s infinite; }
.status-dot.completed { background: var(--color-success); }
.status-dot.error { background: var(--color-error); }

.status-icon {
  font-size: 12px;
}

.status-name {
  font-weight: 510;
  color: var(--color-text-secondary);
}

.status-time {
  color: var(--color-text-subtle);
}

.status-idle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-micro);
  color: var(--color-text-subtle);
}

.idle { background: var(--color-text-subtle); }

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

@media (max-width: 768px) {
  .statusbar {
    display: none;
  }
}
</style>