<template>
  <header class="topbar">
    <button class="hamburger" @click="$emit('toggle-sidebar')" :title="collapsed ? 'Expandir' : 'Colapsar'">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
        <line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="12" x2="20" y2="12"/><line x1="4" y1="18" x2="20" y2="18"/>
      </svg>
    </button>

    <div class="breadcrumb">
      <svg class="bc-dot" viewBox="0 0 24 24" fill="currentColor" width="12" height="12">
        <circle cx="12" cy="12" r="4"/>
      </svg>
      <span class="bc-path">{{ pageTitle }}</span>
    </div>

    <div class="topbar-right">
      <span v-if="weather" class="chip weather">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14">
          <path d="M18 10a6 6 0 0 0-12 0c0 6-6 8-6 8h24s-6-2-6-8z"/>
        </svg>
        {{ weather.temp }}°
      </span>
      <span v-if="exchange" class="chip exchange">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14">
          <line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
        </svg>
        ₡{{ exchange.usdToCrc?.toFixed(0) }}
      </span>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

defineProps({ collapsed: Boolean })
defineEmits(['toggle-sidebar'])

const route = useRoute()
const weather = ref(null)
const exchange = ref(null)

const pageTitle = computed(() => {
  const titles = {
    '/': 'Dashboard',
    '/agents': 'Agentes',
    '/projects': 'Proyectos',
    '/analytics': 'Analytics',
    '/system': 'Sistema',
    '/exchange': 'Exchange',
    '/news': 'Noticias',
  }
  return titles[route.path] || route.meta?.label || route.name || 'Sirius'
})

onMounted(async () => {
  try { weather.value = await fetch('/api/weather').then(r => r.json()) } catch {}
  try { exchange.value = await fetch('/api/exchange').then(r => r.json()) } catch {}
})
</script>

<style scoped>
.topbar {
  height: var(--topbar-height);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px;
  border-bottom: 1px solid var(--border);
  background: var(--canvas);
  flex-shrink: 0;
}

.hamburger {
  display: flex; align-items: center; justify-content: center;
  width: 32px; height: 32px;
  border-radius: var(--radius-md);
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
  flex-shrink: 0;
}
.hamburger:hover { background: var(--hover); color: var(--text-primary); }

.breadcrumb {
  display: flex; align-items: center; gap: 8px;
  flex: 1;
}
.bc-dot { color: var(--brand); flex-shrink: 0; }
.bc-path {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--text-primary);
}

.topbar-right {
  display: flex; gap: 10px;
}

.chip {
  display: flex; align-items: center; gap: 5px;
  font-size: 13px; font-weight: 500;
  padding: 4px 10px;
  border-radius: var(--radius-full);
  background: var(--surface);
  border: 1px solid var(--border);
  color: var(--text-secondary);
}
.chip svg { opacity: 0.6; }
.weather svg { color: var(--warning); }
.exchange svg { color: var(--success); }
</style>
