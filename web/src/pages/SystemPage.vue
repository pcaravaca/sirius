<template>
  <div class="sys-page">
    <div class="sys-grid">
      <div class="sys-card"><h4>⚡ CPU</h4><div class="big">{{ sysInfo.cpu?.used || '?' }}%</div><div class="sub">{{ sysInfo.cpu?.cores || '?' }} cores</div></div>
      <div class="sys-card"><h4>🧠 RAM</h4><div class="big">{{ ramVal }}</div><div class="sub">de {{ totalRam }} GB</div></div>
      <div class="sys-card" v-if="sysInfo.gpu?.loaded"><h4>🎮 GPU</h4><div class="big">{{ gpuVal }}</div><div class="sub">{{ sysInfo.gpu.processes || 0 }} procesos</div></div>
      <div class="sys-card"><h4>💾 Disco</h4><div class="big">{{ diskVal }}</div><div class="sub">de {{ totalDisk }} GB</div></div>
      <div class="sys-card"><h4>⏱ Uptime</h4><div class="big">{{ sysInfo.uptime || '?' }}</div></div>
    </div>
    <div class="ollama-section">
      <h4>🦙 Ollama — {{ sysInfo.ollama?.status || 'desconectado' }}</h4>
      <div v-if="models.length" class="model-list">
        <div v-for="m in models" :key="m.name" class="model-item">{{ m.name }} <span class="model-size">{{ m.size }}</span></div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
const sysInfo = ref({})
const models = ref([])
const ramVal = computed(() => sysInfo.value.ram?.used ? sysInfo.value.ram.used.toFixed(1) + ' GB' : '?')
const totalRam = computed(() => sysInfo.value.ram?.total?.toFixed(0) || '?')
const diskVal = computed(() => sysInfo.value.disk?.used?.toFixed(0) || '?')
const totalDisk = computed(() => sysInfo.value.disk?.total?.toFixed(0) || '?')
const gpuVal = computed(() => sysInfo.value.gpu?.loaded ? sysInfo.value.gpu.vramUsed?.toFixed(1) + ' GB' : '—')
onMounted(async () => {
  try { const s = await fetch('/api/system').then(r => r.json()); sysInfo.value = s; models.value = s.ollama?.models || [] } catch {}
})
</script>
<style scoped>
.sys-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; margin-bottom: 24px; }
.sys-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; }
.sys-card h4 { font-size: 12px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; margin-bottom: 8px; }
.big { font-size: 28px; font-weight: 700; font-family: var(--font-mono); }
.sub { font-size: 13px; color: var(--text-muted); margin-top: 4px; }
.ollama-section { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; }
.ollama-section h4 { margin-bottom: 12px; }
.model-list { display: flex; flex-direction: column; gap: 4px; }
.model-item { display: flex; justify-content: space-between; padding: 6px 10px; background: var(--surface); border-radius: 4px; font-family: var(--font-mono); font-size: 13px; }
.model-size { color: var(--text-muted); }
</style>
