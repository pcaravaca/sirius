<template>
  <header class="topbar">
    <div class="topbar-left">
      <span class="page-title">{{ pageTitle }}</span>
    </div>
    <div class="topbar-right">
      <div class="info-chip" v-if="weather" :title="`${weather.condition}, ${weather.location}`">
        <span class="info-icon">{{ weather.icon || '🌤️' }}</span>
        <span class="info-value">{{ weather.temp }}°</span>
      </div>
      <div class="info-chip" v-if="exchange" title="USD/CRC">
        <span class="info-icon">💵</span>
        <span class="info-value">₡{{ exchange.usdToCrc?.toFixed(0) }}</span>
      </div>
      <div class="ws-indicator" :class="{ connected: wsConnected }" :title="wsConnected ? 'Connected' : 'Disconnected'">
        <span class="ws-dot"></span>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useWeather } from '@/composables/useApi'
import { useExchange } from '@/composables/useApi'
import { useWebSocket } from '@/composables/useWebSocket'

const route = useRoute()
const { data: weather } = useWeather()
const { exchange } = useExchange()
const { isConnected: wsConnected } = useWebSocket()

const pageTitle = computed(() => {
  return route.meta?.label || route.name || 'Sirius'
})

// Fetch weather and exchange on mount
import { onMounted } from 'vue'
const { fetch: fetchWeather } = useWeather()
const { fetch: fetchExchange } = useExchange()

onMounted(() => {
  fetchWeather()
  fetchExchange()
  // Refresh every 5 min
  setInterval(fetchWeather, 300000)
  setInterval(fetchExchange, 300000)
})
</script>

<style scoped>
.topbar {
  height: var(--topbar-height);
  background: var(--color-panel);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-4);
  flex-shrink: 0;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.page-title {
  font-size: var(--font-small);
  font-weight: 590;
  letter-spacing: -0.165px;
  color: var(--color-text-primary);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.info-chip {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  padding: 2px 8px;
  background: var(--color-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-full);
  font-size: var(--font-micro);
  color: var(--color-text-secondary);
}

.info-icon {
  font-size: 12px;
}

.info-value {
  font-weight: 510;
  font-variant-numeric: tabular-nums;
}

.ws-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-error);
  transition: background var(--transition-normal);
}

.ws-indicator.connected {
  background: var(--color-success);
  box-shadow: 0 0 6px var(--color-success);
}

@media (max-width: 768px) {
  .topbar {
    padding: 0 var(--space-3);
  }
  .info-value {
    display: none;
  }
}
</style>