<template>
  <footer class="statusbar">
    <span class="ws-dot" :class="wsConnected ? 'live' : 'dead'"></span>
    <span v-if="activeAgent" class="active-info">
      <strong>{{ activeAgent.name }}</strong> · {{ activeAgent.status }} · {{ activeAgent.task }}
    </span>
    <span class="clock">{{ time }}</span>
  </footer>
</template>
<script setup>
import { ref, inject, onMounted, onUnmounted } from 'vue'
const wsConnected = inject('wsConnected', ref(false))
const activeAgent = ref(null)
const time = ref('')
let ci
onMounted(async () => {
  const update = () => { const n = new Date(); time.value = n.toLocaleTimeString('es-CR', { hour12: false, hour: '2-digit', minute: '2-digit' }) }
  update(); ci = setInterval(update, 30000)
  try { const a = await fetch('/api/agents').then(r => r.json()); activeAgent.value = a.find(x => x.status === 'working') || null } catch {}
})
onUnmounted(() => clearInterval(ci))
</script>
<style scoped>
.statusbar { height: 36px; display: flex; align-items: center; gap: 12px; padding: 0 16px; background: var(--panel); border-top: 1px solid var(--border); font-size: 12px; color: var(--text-muted); flex-shrink: 0; }
.ws-dot { width: 8px; height: 8px; border-radius: 50%; }
.ws-dot.live { background: var(--success); }
.ws-dot.dead { background: var(--text-muted); }
.active-info { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.clock { font-family: var(--font-mono); margin-left: auto; }
</style>
