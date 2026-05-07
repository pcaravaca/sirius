<template>
  <div class="dashboard">
    <!-- Stats Row -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="22" height="22">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-val">{{ projectsCount }}</div>
          <div class="stat-label">Proyectos</div>
        </div>
      </div>
      <div class="stat-card accent">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="22" height="22">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-val">{{ agentsActive }}</div>
          <div class="stat-label">Agentes activos</div>
        </div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="22" height="22">
            <line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-val">{{ opTasks }}</div>
          <div class="stat-label">Tareas OP</div>
        </div>
      </div>
      <div class="stat-card" :class="healthClass">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="22" height="22">
            <rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-val">{{ healthLabel }}</div>
          <div class="stat-label">Sistema</div>
        </div>
      </div>
    </div>

    <!-- System Monitor Gauges -->
    <div class="section-label">Sistema</div>
    <div class="sys-monitor">
      <div class="sys-gauge">
        <div class="gauge-top">
          <span class="gauge-label">CPU</span>
          <span class="gauge-val">{{ sysInfo.cpu?.used || 0 }}%</span>
        </div>
        <div class="gauge-bar">
          <div class="gauge-fill" :style="{ width: (sysInfo.cpu?.used || 0) + '%', background: cpuColor }"></div>
        </div>
        <div class="gauge-sub">{{ sysInfo.cpu?.cores || 0 }} cores</div>
      </div>
      <div class="sys-gauge">
        <div class="gauge-top">
          <span class="gauge-label">RAM</span>
          <span class="gauge-val">{{ ramVal }}</span>
        </div>
        <div class="gauge-bar">
          <div class="gauge-fill" :style="{ width: ramPct + '%', background: ramColor }"></div>
        </div>
        <div class="gauge-sub">de {{ totalRam }} GB</div>
      </div>
      <div v-if="sysInfo.gpu?.loaded" class="sys-gauge">
        <div class="gauge-top">
          <span class="gauge-label">GPU</span>
          <span class="gauge-val">{{ gpuVal }}</span>
        </div>
        <div class="gauge-bar">
          <div class="gauge-fill" :style="{ width: gpuPct + '%', background: gpuColor }"></div>
        </div>
        <div class="gauge-sub">{{ sysInfo.gpu.processes || 0 }} procesos</div>
      </div>
      <div class="sys-gauge">
        <div class="gauge-top">
          <span class="gauge-label">Disco</span>
          <span class="gauge-val">{{ diskVal }}</span>
        </div>
        <div class="gauge-bar">
          <div class="gauge-fill" :style="{ width: diskPct + '%', background: diskColor }"></div>
        </div>
        <div class="gauge-sub">de {{ totalDisk }} GB</div>
      </div>
    </div>

    <!-- Working Agents Panel -->
    <div class="section-label">Agentes trabajando</div>
    <div v-if="workingAgents.length > 0" class="working-panels">
      <div v-for="agent in workingAgents" :key="agent.id" class="working-panel">
        <div class="wp-header">
          <div class="wp-agent-info">
            <div class="wp-name">{{ agent.name }}</div>
            <div class="wp-model">{{ agent.model?.split('/').pop() || '—' }}</div>
          </div>
          <span class="wp-badge">ACTIVO</span>
        </div>
        <div v-if="agent.task" class="wp-task">{{ agent.task }}</div>
        <div v-if="agent.output" class="wp-output">{{ agent.output }}</div>
      </div>
    </div>
    <div v-else class="empty-state">
      <span class="empty-icon">○</span>
      <span class="empty-text">Sin agentes activos</span>
    </div>

    <!-- All Agents Grid -->
    <div class="section-label">Todos los agentes</div>
    <div class="agents-grid">
      <div v-for="a in agents" :key="a.id" class="agent-card" :class="a.status" @click="$router.push('/agents')">
        <div class="ac-top">
          <div class="ac-name">{{ a.name }}</div>
          <div class="ac-dot" :class="a.status"></div>
        </div>
        <div class="ac-model">{{ a.model?.split('/').pop() || '—' }}</div>
        <div v-if="a.task" class="ac-task">{{ a.task }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

// ── State ──
const agents = ref([])
const projectsCount = ref(0)
const opTasks = ref(0)
const sysInfo = ref({})
let pollTimer, sysTimer

// ── Computed ──
const agentsActive = computed(() => agents.value.filter(a => a.status === 'working').length)
const workingAgents = computed(() => agents.value.filter(a => a.status === 'working'))

const healthLabel = computed(() => {
  const cpu = sysInfo.value.cpu?.used || 0
  if (cpu > 80) return 'CRÍTICO'
  if (cpu > 60) return 'ALERTA'
  return 'OK'
})
const healthClass = computed(() => {
  const cpu = sysInfo.value.cpu?.used || 0
  if (cpu > 80) return 'error'
  if (cpu > 60) return 'warning'
  return 'ok'
})

const ramVal = computed(() => sysInfo.value.ram?.used ? sysInfo.value.ram.used.toFixed(1) + ' GB' : '—')
const totalRam = computed(() => {
  const t = parseFloat(sysInfo.value.ram?.total || 0)
  return t > 0 ? t.toFixed(1) : '—'
})
const ramPct = computed(() => {
  const t = parseFloat(sysInfo.value.ram?.total || 1)
  const u = parseFloat(sysInfo.value.ram?.used || 0)
  return t > 0 ? Math.round((u / t) * 100) : 0
})

const diskVal = computed(() => sysInfo.value.disk?.used ? sysInfo.value.disk.used.toFixed(0) + ' GB' : '—')
const totalDisk = computed(() => {
  const t = parseFloat(sysInfo.value.disk?.total || 0)
  return t > 0 ? t.toFixed(0) : '—'
})
const diskPct = computed(() => {
  const t = parseFloat(sysInfo.value.disk?.total || 1)
  const u = parseFloat(sysInfo.value.disk?.used || 0)
  return t > 0 ? Math.round((u / t) * 100) : 0
})

const gpuVal = computed(() => {
  if (!sysInfo.value.gpu?.loaded) return '—'
  return (sysInfo.value.gpu.vramUsed || 0).toFixed(1) + ' GB'
})
const gpuPct = computed(() => {
  if (!sysInfo.value.gpu?.loaded || !sysInfo.value.gpu.vramTotal) return 0
  const t = parseFloat(sysInfo.value.gpu.vramTotal || 1)
  const u = parseFloat(sysInfo.value.gpu.vramUsed || 0)
  return t > 0 ? Math.round((u / t) * 100) : 0
})

// ── Gauge colors ──
const gaugeColor = (pct) => {
  if (pct > 80) return 'var(--error)'
  if (pct > 60) return 'var(--warning)'
  return 'var(--brand)'
}
const cpuColor = computed(() => gaugeColor(sysInfo.value.cpu?.used || 0))
const ramColor = computed(() => gaugeColor(ramPct.value))
const diskColor = computed(() => gaugeColor(diskPct.value))
const gpuColor = computed(() => gaugeColor(gpuPct.value))

// ── API fetches ──
async function fetchAgents() {
  try {
    const res = await fetch('/api/agents')
    const data = await res.json()
    agents.value = Array.isArray(data) ? data : []
  } catch { agents.value = [] }
}

async function fetchProjects() {
  try {
    const res = await fetch('/api/workspace-projects')
    const data = await res.json()
    projectsCount.value = Array.isArray(data) ? data.length : 0
  } catch { projectsCount.value = 0 }
}

async function fetchSystem() {
  try {
    const res = await fetch('/api/system/quick')
    sysInfo.value = await res.json()
  } catch {}
}

async function fetchOpTasks() {
  try {
    const res = await fetch('/api/op-tasks')
    const data = await res.json()
    opTasks.value = Array.isArray(data) ? data.length : 0
  } catch { opTasks.value = 0 }
}

// ── Lifecycle ──
onMounted(async () => {
  await Promise.allSettled([fetchAgents(), fetchProjects(), fetchSystem(), fetchOpTasks()])
  pollTimer = setInterval(fetchAgents, 3000)
  sysTimer = setInterval(fetchSystem, 3000)
})

onUnmounted(() => {
  clearInterval(pollTimer)
  clearInterval(sysTimer)
})
</script>

<style scoped>
.dashboard {
  max-width: 100%;
  overflow-x: hidden;
}

/* ── Stats Row ── */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 20px;
}

.stat-card {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 18px;
  display: flex;
  align-items: center;
  gap: 14px;
  transition: border-color var(--transition-fast);
}
.stat-card:hover { border-color: var(--border-standard); }
.stat-card.ok { border-left: 3px solid var(--success); }
.stat-card.warning { border-left: 3px solid var(--warning); }
.stat-card.error { border-left: 3px solid var(--error); }
.stat-card.accent { border-left: 3px solid var(--brand); }

.stat-icon {
  display: flex; align-items: center; justify-content: center;
  width: 44px; height: 44px;
  background: var(--surface);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  flex-shrink: 0;
}
.stat-card.accent .stat-icon { color: var(--brand); }
.stat-card.warning .stat-icon { color: var(--warning); }
.stat-card.ok .stat-icon { color: var(--success); }
.stat-card.error .stat-icon { color: var(--error); }

.stat-content {
  flex: 1;
  min-width: 0;
}
.stat-val {
  font-size: 32px;
  font-weight: 700;
  font-family: var(--font-mono);
  color: var(--text-primary);
  line-height: 1.1;
}
.stat-label {
  font-size: var(--font-size-xs);
  color: var(--text-muted);
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.5px;
  margin-top: 2px;
}

/* ── Section Label ── */
.section-label {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.6px;
  margin-bottom: 8px;
  margin-top: 4px;
}

/* ── System Monitor ── */
.sys-monitor {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}
.sys-gauge {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 14px;
}
.gauge-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 8px;
}
.gauge-label {
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.4px;
}
.gauge-val {
  font-size: var(--font-size-lg);
  font-weight: 700;
  font-family: var(--font-mono);
  color: var(--text-primary);
}
.gauge-bar {
  height: 6px;
  background: var(--surface);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 4px;
}
.gauge-fill {
  height: 100%;
  border-radius: inherit;
  transition: width 0.6s ease, background 0.6s ease;
  min-width: 2px;
}
.gauge-sub {
  font-size: var(--font-size-xs);
  color: var(--text-muted);
}

/* ── Empty State ── */
.empty-state {
  display: flex; flex-direction: column; align-items: center; gap: 6px;
  padding: 32px 16px;
  color: var(--text-muted);
  font-size: var(--font-size-sm);
}
.empty-icon { font-size: 24px; opacity: 0.4; }

/* ── Working Panels ── */
.working-panels {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}
.working-panel {
  background: var(--panel);
  border: 1px solid var(--border);
  border-left: 3px solid var(--brand);
  border-radius: var(--radius-md);
  padding: 14px 16px;
}
.wp-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 6px;
}
.wp-name {
  font-size: var(--font-size-lg);
  font-weight: 700;
  color: var(--text-primary);
}
.wp-model {
  font-size: var(--font-size-xs);
  font-family: var(--font-mono);
  color: var(--text-muted);
}
.wp-badge {
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: var(--brand);
  background: rgba(0,153,255,0.1);
  padding: 2px 10px;
  border-radius: var(--radius-full);
  letter-spacing: 0.3px;
}
.wp-task {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-bottom: 6px;
  padding: 8px 10px;
  background: var(--surface);
  border-radius: var(--radius-sm);
  border-left: 2px solid var(--brand);
}
.wp-output {
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  background: #0a0e14;
  border-radius: var(--radius-sm);
  padding: 10px;
  max-height: 120px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.4;
}

/* ── Agents Grid ── */
.agents-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
  margin-bottom: 16px;
}
.agent-card {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 12px 14px;
  cursor: pointer;
  transition: border-color var(--transition-fast), background var(--transition-fast);
  min-width: 0;
}
.agent-card:hover {
  border-color: var(--border-standard);
  background: var(--surface);
}
.agent-card.working { border-left: 3px solid var(--brand); }
.agent-card.error { border-left: 3px solid var(--error); }
.agent-card.completed { border-left: 3px solid var(--success); }

.ac-top {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 2px;
}
.ac-name {
  font-size: var(--font-size-lg);
  font-weight: 700;
  color: var(--text-primary);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.ac-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
}
.ac-dot.idle { background: var(--text-muted); }
.ac-dot.working { background: var(--brand); box-shadow: 0 0 6px var(--brand); }
.ac-dot.error { background: var(--error); box-shadow: 0 0 6px rgba(231,76,60,0.4); }
.ac-dot.completed { background: var(--success); }

.ac-model {
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
  color: var(--text-muted);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.ac-task {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  margin-top: 6px;
  padding-top: 6px;
  border-top: 1px solid var(--border);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}

/* ── Responsive ── */
@media (max-width: 1200px) {
  .stats-row { grid-template-columns: repeat(2, 1fr); }
  .agents-grid { grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 768px) {
  .stats-row { grid-template-columns: 1fr; }
  .agents-grid { grid-template-columns: repeat(2, 1fr); }
  .sys-monitor { grid-template-columns: repeat(2, 1fr); }
}

/* ── CRITICAL: Prevent shift on data update ── */
.stat-val, .gauge-val, .ac-name, .wp-name, .wp-output {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}
.wp-output, .wp-task { white-space: normal; }
</style>
