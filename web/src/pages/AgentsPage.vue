<template>
  <div class="agents-page">
    <div class="filter-bar">
      <button v-for="f in filters" :key="f.key" @click="active=f.key" :class="{active:active===f.key}">{{ f.label }}</button>
    </div>
    <div class="agents-grid">
      <div v-for="a in filtered" :key="a.id" class="agent-card" :class="a.status" @click="selected=a">
        <div class="card-header">
          <strong>{{ a.name }}</strong>
          <span class="status-dot" :class="a.status"></span>
        </div>
        <div class="model">{{ shortModel(a.model) }}</div>
        <div v-if="a.task" class="task">{{ a.task }}</div>
        <div v-if="a.lastOutput" class="output">{{ a.lastOutput.substring(0,100) }}</div>
        <div v-if="a.tokensIn" class="metrics">📊 {{ a.tokensIn }} in · {{ a.tokensOut }} out ⏱ {{ ms(a.durationMs) }}</div>
      </div>
    </div>
    <div v-if="selected" class="modal" @click.self="selected=null">
      <div class="modal-panel">
        <h2>{{ selected.name }}</h2>
        <p>{{ selected.model }} · {{ selected.provider }} · {{ selected.channel }}</p>
        <p><strong>Task:</strong> {{ selected.task || '—' }}</p>
        <pre class="output-box">{{ selected.output || 'No output' }}</pre>
        <button @click="dismiss(selected.id)">Dismiss</button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
const agents = ref([])
const active = ref('all')
const selected = ref(null)
const filters = [{key:'all',label:'Todos'},{key:'working',label:'Activos'},{key:'error',label:'Error'},{key:'completed',label:'Completados'}]
const filtered = computed(() => active.value === 'all' ? agents.value : agents.value.filter(a => a.status === active.value))
function shortModel(m) { return m ? m.split('/').pop() : '—' }
function ms(d) { if (!d) return ''; const s = Math.floor(d/1000); return s < 60 ? `${s}s` : `${Math.floor(s/60)}m ${s%60}s` }
async function dismiss(id) { await fetch('/api/webhook/agent-activity', {method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:id,task:'',status:'idle',output:''})}); selected.value = null; agents.value = await fetch('/api/agents').then(r=>r.json()) }
onMounted(async () => { agents.value = await fetch('/api/agents').then(r=>r.json()) })
</script>
<style scoped>
.filter-bar { display: flex; gap: 4px; margin-bottom: 12px; }
.filter-bar button { padding: 6px 14px; border: 1px solid var(--border); border-radius: 999px; background: transparent; color: var(--text-muted); cursor: pointer; }
.filter-bar button.active { background: var(--brand); color: white; border-color: var(--brand); }
.agents-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 10px; }
.agent-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 14px; cursor: pointer; }
.agent-card:hover { background: var(--surface); }
.agent-card.working { border-left: 3px solid var(--success); }
.agent-card.error { border-left: 3px solid var(--error); }
.agent-card.completed { border-left: 3px solid var(--info); }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.status-dot { width: 10px; height: 10px; border-radius: 50%; }
.status-dot.working { background: var(--success); }
.status-dot.idle { background: var(--text-muted); }
.status-dot.error { background: var(--error); }
.status-dot.completed { background: var(--info); }
.model { font-size: 12px; color: var(--text-muted); font-family: var(--font-mono); margin-top: 4px; }
.task { font-size: 13px; margin-top: 8px; padding: 6px; background: var(--surface); border-radius: 4px; }
.output { font-size: 12px; color: var(--text-muted); font-family: var(--font-mono); margin-top: 4px; }
.metrics { font-size: 11px; color: var(--text-muted); margin-top: 6px; }
.modal { position: fixed; inset: 0; background: rgba(0,0,0,0.6); z-index: 200; display: flex; align-items: center; justify-content: center; }
.modal-panel { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 24px; max-width: 600px; width: 90%; }
.output-box { background: #0a0a0a; padding: 16px; border-radius: var(--radius-md); font-family: var(--font-mono); font-size: 13px; max-height: 300px; overflow-y: auto; white-space: pre-wrap; }
</style>
