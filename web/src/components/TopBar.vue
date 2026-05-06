<template>
  <header class="topbar">
    <h1 class="page-title">{{ pageTitle }}</h1>
    <div class="topbar-right">
      <span v-if="weather" class="chip">🌡️ {{ weather.temp }}°</span>
      <span v-if="exchange" class="chip">₡{{ exchange.usdToCrc?.toFixed(0) }}</span>
    </div>
  </header>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const weather = ref(null)
const exchange = ref(null)
const pageTitle = computed(() => route.meta?.label || route.name || 'Sirius')
onMounted(async () => {
  try { weather.value = await fetch('/api/weather').then(r => r.json()) } catch {}
  try { exchange.value = await fetch('/api/exchange').then(r => r.json()) } catch {}
})
</script>
<style scoped>
.topbar { height: var(--topbar-height); display: flex; align-items: center; justify-content: space-between; padding: 0 24px; border-bottom: 1px solid var(--border); background: var(--canvas); flex-shrink: 0; }
.page-title { font-size: 16px; font-weight: 600; }
.topbar-right { display: flex; gap: 12px; }
.chip { font-size: 13px; color: var(--text-muted); }
</style>
