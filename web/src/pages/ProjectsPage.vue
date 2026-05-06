<template>
  <div class="proj-page">
    <div class="proj-grid">
      <div v-for="p in projects" :key="p.name" class="proj-card">
        <div class="proj-name">{{ p.name }}</div>
        <div class="proj-source">{{ p.source }}</div>
        <div class="proj-lang">{{ p.language }}</div>
        <div v-if="p.gitBranch" class="proj-git">🔀 {{ p.gitBranch }} · {{ (p.lastCommit || '').substring(0,50) }}</div>
        <div class="proj-meta">
          <span>🧪 {{ p.hasTests ? '✅' : '—' }}</span>
          <span>🐳 {{ p.hasDocker ? '✅' : '—' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
const projects = ref([])
onMounted(async () => { projects.value = await fetch('/api/workspace-projects').then(r => r.json()) })
</script>
<style scoped>
.proj-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 12px; }
.proj-card { background: var(--panel); border: 1px solid var(--border); border-left: 3px solid var(--brand); border-radius: var(--radius-md); padding: 16px; }
.proj-name { font-size: 16px; font-weight: 700; }
.proj-source { font-size: 12px; color: var(--text-muted); }
.proj-lang { display: inline-block; font-size: 10px; font-weight: 700; padding: 2px 8px; border-radius: 999px; background: var(--surface); margin-top: 4px; text-transform: uppercase; }
.proj-git { font-size: 12px; font-family: var(--font-mono); color: var(--text-muted); margin-top: 8px; }
.proj-meta { display: flex; gap: 16px; margin-top: 8px; font-size: 12px; color: var(--text-muted); }
</style>
