<template>
  <div class="dashboard">
    <!-- Stats Row -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-val">{{ projects }}</div>
        <div class="stat-label">Proyectos</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{{ agentsActive }}</div>
        <div class="stat-label">Agentes activos</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{{ opTasks }}</div>
        <div class="stat-label">Tareas OP</div>
      </div>
      <div class="stat-card">
        <div class="stat-val" :class="healthClass">{{ healthLabel }}</div>
        <div class="stat-label">Sistema</div>
      </div>
    </div>

    <!-- System Monitor -->
    <div class="section-header">
      <h3 class="section-title">Sistema en vivo</h3>
      <button class="section-link" @click="$router.push('/system')">Ver detalle →</button>
    </div>
    <div class="sys-monitor">
      <div class="sys-gauge">
        <div class="gauge-label"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><rect x="4" y="4" width="16" height="16" rx="2"/><rect x="9" y="9" width="6" height="6"/></svg> CPU</div>
        <div class="gauge-bar"><div class="gauge-fill" :style="{ width: cpuPct + '%', background: cpuColor, boxShadow: '0 0 6px ' + cpuColor }"></div></div>
        <div class="gauge-val">{{ cpuPct }}%</div>
        <div class="gauge-sub">{{ sysInfo.cpu?.cores || '?' }} cores</div>
      </div>
      <div class="sys-gauge">
        <div class="gauge-label"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="M4 20h16"/></svg> RAM</div>
        <div class="gauge-bar"><div class="gauge-fill" :style="{ width: ramPct + '%', background: ramColor }"></div></div>
        <div class="gauge-val">{{ ramVal }}</div>
        <div class="gauge-sub">/ {{ totalRam }} GB · {{ ramPct }}%</div>
      </div>
      <div class="sys-gauge" v-if="gpuLoaded">
        <div class="gauge-label"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg> GPU</div>
        <div class="gauge-bar"><div class="gauge-fill" :style="{ width: gpuPct + '%', background: gpuColor }"></div></div>
        <div class="gauge-val">{{ gpuVal }}</div>
        <div class="gauge-sub">{{ sysInfo.gpu?.processes || 0 }} procesos</div>
      </div>
      <div class="sys-gauge">
        <div class="gauge-label"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><ellipse cx="12" cy="12" rx="10" ry="10"/><circle cx="12" cy="12" r="4"/><circle cx="12" cy="12" r="1"/></svg> Disco</div>
        <div class="gauge-bar"><div class="gauge-fill" :style="{ width: diskPct + '%', background: diskColor }"></div></div>
        <div class="gauge-val">{{ diskVal }}</div>
        <div class="gauge-sub">/ {{ totalDisk }} GB · {{ diskPct }}%</div>
      </div>
      <div class="sys-gauge">
        <div class="gauge-label"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg> Ollama</div>
        <div class="gauge-bar"><div class="gauge-fill" :style="{ width: ollamaOnline ? '100%' : '20%', background: ollamaOnline ? 'var(--success)' : 'var(--text-muted)' }"></div></div>
        <div class="gauge-val">{{ ollamaOnline ? '🟢' : '🔴' }} {{ modelsCount }}</div>
        <div class="gauge-sub">{{ sysInfo.uptime || '' }}</div>
      </div>
    </div>

    <!-- Agents Grid -->
    <div class="section-header">
      <h3 class="section-title">Agentes</h3>
      <button class="section-link" @click="$router.push('/agents')">Ver todos →</button>
    </div>
    <div class="agents-grid" :class="{ 'has-working': workingAgents.length > 0 }">
      <div v-for="a in agents" :key="a.id" class="agent-card" :class="a.status" @click="$router.push('/agents')">
        <div class="agent-name">{{ a.name }}</div>
        <div class="agent-model">{{ shortModel(a.model) }}</div>
        <div class="agent-status-dot" :class="a.status"></div>
        <div class="agent-task" v-if="a.task">{{ a.task }}</div>
      </div>
    </div>

    <!-- Bottom row -->
    <div class="bottom-row">
      <div class="bottom-card">
        <h4>Correos</h4>
        <div v-for="m in emails" :key="m.id" class="email-row">
          <span class="email-from">{{ m.from }}</span>
          <span class="email-subj">{{ m.subject }}</span>
        </div>
      </div>
      <div class="bottom-card">
        <h4>Agenda</h4>
        <div v-for="ev in agenda" :key="ev.summary" class="agenda-row">
          <span class="agenda-date">{{ formatDate(ev.start) }}</span>
          <span class="agenda-text">{{ ev.summary }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const agents = ref([])
const projects = ref(0)
const opTasks = ref(0)
const healthLabel = ref('—')
const healthClass = ref('')
const emails = ref([])
const agenda = ref([])
const sysInfo = ref({})
const modelsCount = ref(0)
let pollTimer = null
let sysTimer = null

function shortModel(m) {
  if (!m) return '—'
  return m.split('/').pop()
}

function formatDate(d) {
  if (!d) return ''
  return d.includes('T') ? d.substring(5, 16) : d.substring(5)
}

const workingAgents = computed(() => agents.value.filter(a => a.status === 'working'))
const agentsActive = computed(() => workingAgents.value.length)

const cpuPct = computed(() => sysInfo.value.cpu?.used || 0)
const ramPct = computed(() => { const r = sysInfo.value.ram; return r?.total ? Math.round((r.used / r.total) * 100) : 0 })
const ramVal = computed(() => { const r = sysInfo.value.ram; return r?.used ? `${r.used.toFixed(1)} GB` : '?' })
const totalRam = computed(() => sysInfo.value.ram?.total?.toFixed(0) || '?')
const diskPct = computed(() => { const d = sysInfo.value.disk; return d?.total ? Math.round((d.used / d.total) * 100) : 0 })
const diskVal = computed(() => { const d = sysInfo.value.disk; return d?.used ? `${d.used.toFixed(0)} GB` : '?' })
const totalDisk = computed(() => sysInfo.value.disk?.total?.toFixed(0) || '?')
const gpuLoaded = computed(() => sysInfo.value.gpu?.loaded)
const gpuPct = computed(() => { const g = sysInfo.value.gpu; return g?.vramTotal ? Math.round((g.vramUsed / g.vramTotal) * 100) : 0 })
const gpuVal = computed(() => { const g = sysInfo.value.gpu; return g?.loaded ? `${g.vramUsed?.toFixed(1)} GB` : '—' })
const ollamaOnline = computed(() => sysInfo.value.ollama?.loaded)

const cpuColor = computed(() => cpuPct.value > 80 ? 'var(--error)' : cpuPct.value > 50 ? 'var(--warning)' : 'var(--success)')
const ramColor = computed(() => ramPct.value > 80 ? 'var(--error)' : ramPct.value > 50 ? 'var(--warning)' : 'var(--success)')
const gpuColor = computed(() => gpuPct.value > 80 ? 'var(--error)' : gpuPct.value > 50 ? 'var(--warning)' : 'var(--success)')
const diskColor = computed(() => diskPct.value > 80 ? 'var(--error)' : diskPct.value > 50 ? 'var(--warning)' : 'var(--success)')

async function fetchAll() {
  try {
    agents.value = await fetch('/api/agents').then(r => r.json())
  } catch {}
  try {
    const p = await fetch('/api/projects').then(r => r.json())
    projects.value = Array.isArray(p) ? p.length : 0
  } catch {}
  try {
    const t = await fetch('/api/op-tasks').then(r => r.json())
    opTasks.value = Array.isArray(t) ? t.length : 0
  } catch {}
  try {
    const e = await fetch('/api/emails').then(r => r.json())
    emails.value = Array.isArray(e) ? e.slice(0, 5) : []
  } catch {}
  try {
    const a = await fetch('/api/agenda').then(r => r.json())
    agenda.value = Array.isArray(a) ? a.slice(0, 8) : []
  } catch {}
}

async function loadSystem() {
  try {
    const sys = await fetch('/api/system/quick').then(r => r.json())
    sysInfo.value = sys
    if (sys.ollama?.loaded) { healthLabel.value = 'OK'; healthClass.value = 'ok'; modelsCount.value = sys.ollama.models?.length || 0 }
    else { healthLabel.value = '—'; modelsCount.value = 0 }
  } catch { healthLabel.value = '—' }
}

onMounted(async () => {
  await Promise.all([fetchAll(), loadSystem()])
  pollTimer = setInterval(fetchAll, 3000)
  sysTimer = setInterval(loadSystem, 3000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
  if (sysTimer) clearInterval(sysTimer)
})
</script>

<style scoped>
.dashboard { max-width: 1200px; }

.stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 24px; }
.stat-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; }
.stat-val { font-size: 28px; font-weight: 700; color: var(--text-primary); font-family: var(--font-mono); }
.stat-val.ok { color: var(--success); }
.stat-label { font-size: 12px; color: var(--text-muted); margin-top: 4px; }

.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.section-title { font-size: 18px; font-weight: 700; }
.section-link { background: none; border: none; color: var(--brand-interactive); font-size: 13px; cursor: pointer; padding: 4px 8px; border-radius: var(--radius-md); }
.section-link:hover { background: var(--hover); }

.sys-monitor { display: grid; grid-template-columns: repeat(5, 1fr); gap: 10px; margin-bottom: 24px; }
.sys-gauge { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 14px; }
.sys-gauge:first-child { background: linear-gradient(135deg, var(--panel) 0%, rgba(94,106,210,0.08) 100%); border-color: rgba(94,106,210,0.3); }
.gauge-label { font-size: 11px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 6px; display: flex; align-items: center; gap: 4px; }
.gauge-bar { height: 8px; background: var(--surface); border-radius: 4px; overflow: hidden; margin-bottom: 6px; }
.gauge-fill { height: 100%; border-radius: 4px; transition: width 0.5s ease, box-shadow 0.5s ease; }
.gauge-val { font-size: 20px; font-weight: 700; font-family: var(--font-mono); line-height: 1.2; }
.gauge-sub { font-size: 11px; color: var(--text-muted); display: block; margin-top: 2px; }

.agents-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; margin-bottom: 24px; }
.agents-grid.has-working { opacity: 0.6; }
.agent-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 12px; cursor: pointer; transition: background var(--transition-fast); }
.agent-card:hover { background: var(--surface); }
.agent-card.working { border-left: 3px solid var(--success); }
.agent-name { font-weight: 700; font-size: 14px; }
.agent-model { font-size: 11px; color: var(--text-muted); font-family: var(--font-mono); }
.agent-status-dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.agent-status-dot.working { background: var(--success); }
.agent-status-dot.idle { background: var(--text-muted); }
.agent-task { font-size: 12px; margin-top: 4px; color: var(--text-secondary); }

.bottom-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.bottom-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; }
.bottom-card h4 { font-size: 13px; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; margin-bottom: 12px; }
.email-row, .agenda-row { font-size: 13px; padding: 4px 0; border-bottom: 1px solid var(--border-subtle); display: flex; gap: 8px; }
.email-from { font-weight: 600; min-width: 80px; }
.email-subj, .agenda-text { color: var(--text-muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.agenda-date { font-family: var(--font-mono); font-size: 11px; color: var(--text-muted); }
</style>
