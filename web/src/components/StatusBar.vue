<template>
  <footer class="statusbar">
    <div class="sb-left">
      <span class="ws-dot" :class="wsConnected ? 'live' : 'dead'"></span>
      <span class="ws-label">{{ wsConnected ? 'Conectado' : 'Desconectado' }}</span>
    </div>

    <div class="sb-center" v-if="activeAgent">
      <span class="active-pill">{{ activeAgent.name }}</span>
      <span class="active-status">{{ getStatusLabel(activeAgent.status) }}</span>
      <span v-if="activeAgent.task" class="active-task">{{ activeAgent.task }}</span>
    </div>

    <div class="sb-right">
      <span class="clock">{{ time }}</span>
    </div>
  </footer>
</template>

<script setup>
import { ref, inject, onMounted, onUnmounted } from 'vue'

const wsConnected = inject('wsConnected', ref(false))
const activeAgent = ref(null)
const time = ref('')
let ci

function getStatusLabel(status) {
  const labels = { working: 'ACTIVO', idle: 'INACTIVO', completed: 'HECHO', error: 'ERROR' }
  return labels[status] || status?.toUpperCase() || ''
}

onMounted(async () => {
  const update = () => {
    const n = new Date()
    time.value = n.toLocaleTimeString('es-CR', { hour12: false, hour: '2-digit', minute: '2-digit' })
  }
  update(); ci = setInterval(update, 30000)
  try {
    const res = await fetch('/api/agents')
    const data = await res.json()
    activeAgent.value = (Array.isArray(data) ? data : []).find(a => a.status === 'working') || null
  } catch {}
})

onUnmounted(() => clearInterval(ci))
</script>

<style scoped>
.statusbar {
  height: var(--statusbar-height);
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 16px;
  background: #0a0e14;
  border-top: 1px solid var(--border);
  font-size: var(--font-size-sm);
  flex-shrink: 0;
}

.sb-left {
  display: flex; align-items: center; gap: 6px;
}
.ws-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
}
.ws-dot.live { background: var(--success); box-shadow: 0 0 6px rgba(0,168,107,0.4); }
.ws-dot.dead { background: var(--text-muted); }
.ws-label { color: var(--text-muted); font-size: var(--font-size-xs); }

.sb-center {
  flex: 1;
  display: flex; align-items: center; gap: 8px;
  overflow: hidden;
}
.active-pill {
  background: rgba(0,168,107,0.12);
  color: var(--success);
  padding: 1px 8px;
  border-radius: var(--radius-full);
  font-weight: 700;
  font-size: var(--font-size-xs);
  flex-shrink: 0;
}
.active-status {
  color: var(--text-muted);
  font-size: var(--font-size-xs);
  flex-shrink: 0;
}
.active-task {
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sb-right { margin-left: auto; }
.clock {
  font-family: var(--font-mono);
  color: var(--text-muted);
  font-size: var(--font-size-sm);
}
</style>
