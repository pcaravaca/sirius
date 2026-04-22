<template>
  <div class="exchange-page">
    <div class="page-header">
      <h2>💱 Exchange Rate</h2>
      <span class="badge" v-if="loading">Loading...</span>
    </div>

    <div class="rate-cards" v-if="exchange">
      <div class="rate-card main">
        <span class="rate-label">USD → CRC</span>
        <span class="rate-value">₡{{ exchange.usdToCrc?.toFixed(2) }}</span>
        <span class="rate-sub">Buy: ₡{{ exchange.usdToCrcBuy?.toFixed(2) || '—' }}</span>
      </div>
      <div class="rate-card">
        <span class="rate-label">Source</span>
        <span class="rate-value small">{{ exchange.source || 'BCCR' }}</span>
      </div>
      <div class="rate-card">
        <span class="rate-label">Updated</span>
        <span class="rate-value small">{{ formatTime(exchange.updatedAt) }}</span>
      </div>
    </div>

    <div class="chart-section" v-if="history && history.length">
      <h3>30 Day History</h3>
      <div class="mini-chart">
        <div
          v-for="(item, i) in history"
          :key="i"
          class="chart-bar"
          :style="{ height: barHeight(item.rate) + '%' }"
          :title="`${item.date}: ₡${item.rate?.toFixed(2)}`"
        ></div>
      </div>
    </div>

    <div class="predictions" v-if="predictions && predictions.length">
      <h3>Predictions</h3>
      <div class="pred-list">
        <div v-for="(p, i) in predictions.slice(0, 7)" :key="i" class="pred-item">
          <span class="pred-date">{{ p.date }}</span>
          <span class="pred-rate">₡{{ p.rate?.toFixed(2) }}</span>
          <span class="pred-dir" :class="p.direction">{{ p.direction === 'up' ? '↑' : '↓' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useExchange } from '@/composables/useApi'

const { exchange, history, predictions, loading, fetch } = useExchange()

function barHeight(rate) {
  if (!rate || !history?.length) return 0
  const rates = history.map(h => h.rate).filter(Boolean)
  const min = Math.min(...rates)
  const max = Math.max(...rates)
  const range = max - min || 1
  return Math.max(5, ((rate - min) / range) * 100)
}

function formatTime(ts) {
  if (!ts) return '—'
  return new Date(ts).toLocaleTimeString('es-CR', { hour: '2-digit', minute: '2-digit' })
}

onMounted(fetch)
</script>

<style scoped>
.exchange-page {
  max-width: 900px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-6);
}

.page-header h2 {
  font-size: var(--font-h3);
  font-weight: 590;
}

.badge {
  font-size: var(--font-micro);
  color: var(--color-text-muted);
}

.rate-cards {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: var(--space-4);
  margin-bottom: var(--space-8);
}

.rate-card {
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.rate-card.main {
  background: rgba(94, 106, 210, 0.08);
  border-color: rgba(94, 106, 210, 0.2);
}

.rate-label {
  font-size: var(--font-micro);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.rate-value {
  font-size: 32px;
  font-weight: 590;
  color: var(--color-text-primary);
  font-variant-numeric: tabular-nums;
}

.rate-value.small {
  font-size: var(--font-body);
}

.rate-sub {
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.chart-section, .predictions {
  margin-bottom: var(--space-8);
}

.chart-section h3, .predictions h3 {
  font-size: var(--font-small);
  font-weight: 510;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-4);
}

.mini-chart {
  display: flex;
  align-items: flex-end;
  gap: 3px;
  height: 120px;
  padding: var(--space-3);
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
}

.chart-bar {
  flex: 1;
  background: var(--color-accent);
  border-radius: 2px 2px 0 0;
  min-height: 4px;
  opacity: 0.7;
  transition: opacity var(--transition-fast);
}

.chart-bar:hover {
  opacity: 1;
}

.pred-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.pred-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  background: var(--color-surface);
  border-radius: var(--radius-md);
  font-size: var(--font-caption);
}

.pred-date {
  color: var(--color-text-muted);
  width: 100px;
}

.pred-rate {
  flex: 1;
  color: var(--color-text-primary);
  font-variant-numeric: tabular-nums;
  font-family: var(--font-mono);
}

.pred-dir.up { color: var(--color-success); }
.pred-dir.down { color: var(--color-error); }

@media (max-width: 768px) {
  .rate-cards {
    grid-template-columns: 1fr;
  }
}
</style>