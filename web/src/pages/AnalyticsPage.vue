<script setup>
import { ref, computed } from 'vue'

// Placeholder data until analytics API is connected
const tokenUsage = ref([
  { label: 'Mon', value: 12400, max: 20000 },
  { label: 'Tue', value: 8700, max: 20000 },
  { label: 'Wed', value: 18200, max: 20000 },
  { label: 'Thu', value: 15100, max: 20000 },
  { label: 'Fri', value: 9300, max: 20000 },
  { label: 'Sat', value: 4200, max: 20000 },
  { label: 'Sun', value: 3100, max: 20000 },
])

const costBreakdown = ref([
  { label: 'LLM Inference', value: 74, color: 'var(--color-accent)' },
  { label: 'Embeddings', value: 14, color: 'var(--color-info)' },
  { label: 'Storage', value: 8, color: 'var(--color-warning)' },
  { label: 'Network', value: 4, color: 'var(--color-success)' },
])

const taskCompletion = ref([
  { label: 'Completed', value: 68, color: 'var(--color-success)' },
  { label: 'In Progress', value: 18, color: 'var(--color-accent)' },
  { label: 'Blocked', value: 6, color: 'var(--color-error)' },
  { label: 'Todo', value: 8, color: 'var(--color-text-muted)' },
])

const totalTokens = computed(() =>
  tokenUsage.value.reduce((sum, d) => sum + d.value, 0)
)

const estCost = ref(24.50)

function barHeight(item) {
  return `${Math.round((item.value / item.max) * 100)}%`
}

function formatTokens(n) {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`
  return n
}
</script>

<template>
  <div class="analytics-page">
    <header class="page-header">
      <h1 class="page-title">Analytics</h1>
      <span class="page-badge">Preview</span>
    </header>

    <!-- Summary row -->
    <section class="summary-row">
      <div class="summary-card">
        <span class="summary-value">{{ formatTokens(totalTokens) }}</span>
        <span class="summary-label">Tokens this week</span>
      </div>
      <div class="summary-card">
        <span class="summary-value">${{ estCost.toFixed(2) }}</span>
        <span class="summary-label">Est. cost</span>
      </div>
      <div class="summary-card">
        <span class="summary-value">68%</span>
        <span class="summary-label">Task completion</span>
      </div>
    </section>

    <!-- Token Usage Bar Chart -->
    <section class="panel">
      <h2 class="panel-title">Token Usage</h2>
      <span class="panel-subtitle">Daily token consumption this week</span>

      <div class="bar-chart">
        <div class="bar-col" v-for="item in tokenUsage" :key="item.label">
          <div class="bar-wrapper">
            <div
              class="bar-fill"
              :style="{ height: barHeight(item) }"
            ></div>
          </div>
          <span class="bar-label">{{ item.label }}</span>
          <span class="bar-value">{{ formatTokens(item.value) }}</span>
        </div>
      </div>
    </section>

    <!-- Two column: Cost + Completion -->
    <div class="two-col">
      <!-- Cost Breakdown -->
      <section class="panel">
        <h2 class="panel-title">Cost Breakdown</h2>
        <span class="panel-subtitle">Estimated spending distribution</span>

        <div class="h-bar-list">
          <div class="h-bar-item" v-for="item in costBreakdown" :key="item.label">
            <div class="h-bar-header">
              <span class="h-bar-label">{{ item.label }}</span>
              <span class="h-bar-value">{{ item.value }}%</span>
            </div>
            <div class="h-bar-track">
              <div
                class="h-bar-fill"
                :style="{
                  width: item.value + '%',
                  background: item.color,
                }"
              ></div>
            </div>
          </div>
        </div>
      </section>

      <!-- Task Completion Rate -->
      <section class="panel">
        <h2 class="panel-title">Task Completion</h2>
        <span class="panel-subtitle">Task status distribution</span>

        <div class="h-bar-list">
          <div class="h-bar-item" v-for="item in taskCompletion" :key="item.label">
            <div class="h-bar-header">
              <span class="h-bar-label">{{ item.label }}</span>
              <span class="h-bar-value">{{ item.value }}%</span>
            </div>
            <div class="h-bar-track">
              <div
                class="h-bar-fill"
                :style="{
                  width: item.value + '%',
                  background: item.color,
                }"
              ></div>
            </div>
          </div>
        </div>

        <!-- Stacked bar -->
        <div class="stacked-bar">
          <div
            class="stacked-segment"
            v-for="item in taskCompletion"
            :key="item.label"
            :style="{
              width: item.value + '%',
              background: item.color,
            }"
          ></div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.analytics-page {
  padding: var(--space-8);
  max-width: 1120px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-8);
}

.page-title {
  font-family: var(--font-primary);
  font-size: var(--font-h2);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.page-badge {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  font-weight: 500;
  color: var(--color-accent);
  background: rgba(94, 106, 210, 0.1);
  border: 1px solid rgba(94, 106, 210, 0.2);
  border-radius: var(--radius-full);
  padding: 2px var(--space-2);
}

/* Summary row */
.summary-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.summary-card {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.summary-value {
  font-family: var(--font-primary);
  font-size: var(--font-h2);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.summary-label {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

/* Panel */
.panel {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  margin-bottom: var(--space-6);
}

.panel-title {
  font-family: var(--font-primary);
  font-size: var(--font-small);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-1) 0;
}

.panel-subtitle {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
  display: block;
  margin-bottom: var(--space-5);
}

/* Vertical bar chart */
.bar-chart {
  display: flex;
  align-items: stretch;
  gap: var(--space-4);
  height: 200px;
  padding-top: var(--space-3);
}

.bar-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
}

.bar-wrapper {
  flex: 1;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  background: var(--color-surface);
  border-radius: var(--radius-md);
  overflow: hidden;
  position: relative;
}

.bar-fill {
  background: var(--color-accent);
  border-radius: var(--radius-md);
  transition: height var(--transition-normal);
  min-height: 4px;
}

.bar-label {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  color: var(--color-text-muted);
}

.bar-value {
  font-family: var(--font-mono);
  font-size: var(--font-micro);
  color: var(--color-text-secondary);
}

/* Horizontal bars */
.h-bar-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.h-bar-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.h-bar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.h-bar-label {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
}

.h-bar-value {
  font-family: var(--font-mono);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.h-bar-track {
  height: 8px;
  background: var(--color-surface);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.h-bar-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width var(--transition-normal);
}

/* Stacked bar */
.stacked-bar {
  display: flex;
  height: 8px;
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-top: var(--space-5);
  background: var(--color-surface);
}

.stacked-segment {
  height: 100%;
  transition: width var(--transition-normal);
}

.stacked-segment:first-child {
  border-radius: var(--radius-full) 0 0 var(--radius-full);
}

.stacked-segment:last-child {
  border-radius: 0 var(--radius-full) var(--radius-full) 0;
}

/* Two column layout */
.two-col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-6);
}

.two-col .panel {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .analytics-page {
    padding: var(--space-4);
  }

  .summary-row {
    grid-template-columns: 1fr;
  }

  .two-col {
    grid-template-columns: 1fr;
  }

  .bar-chart {
    height: 150px;
  }
}
</style>