<script setup>
import { ref, onMounted, computed } from 'vue'
import { useSystem } from '@/composables/useApi'

const { system, loading, fetch } = useSystem()
const fetchError = ref(false)

onMounted(async () => {
  try {
    await fetch()
  } catch (e) {
    // Endpoint may 404 — use placeholder data
    fetchError.value = true
  }
})

// Placeholder data when the backend endpoint doesn't exist yet
const placeholder = computed(() => ({
  ollama: {
    loaded: true,
    models: [
      { name: 'llama3:8b', size: '4.7 GB', vram: 48 },
      { name: 'mistral:7b', size: '4.1 GB', vram: 42 },
      { name: 'codellama:13b', size: '7.4 GB', vram: 76 },
    ],
    vram: {
      used: 16.2,
      total: 24,
    },
  },
  system: {
    cpu: { used: 45, cores: 8 },
    ram: { used: 12.4, total: 32 },
    disk: { used: 186, total: 512 },
  },
}))

const data = computed(() => system.value ?? placeholder.value)

const vramPercent = computed(() => {
  const d = data.value?.ollama?.vram
  if (!d) return 0
  return Math.round((d.used / d.total) * 100)
})

const cpuPercent = computed(() => data.value?.system?.cpu?.used ?? 0)
const ramPercent = computed(() => {
  const d = data.value?.system?.ram
  if (!d) return 0
  return Math.round((d.used / d.total) * 100)
})
const diskPercent = computed(() => {
  const d = data.value?.system?.disk
  if (!d) return 0
  return Math.round((d.used / d.total) * 100)
})

function gaugeColor(percent) {
  if (percent >= 90) return 'var(--color-error)'
  if (percent >= 70) return 'var(--color-warning)'
  return 'var(--color-accent)'
}
</script>

<template>
  <div class="system-page">
    <header class="page-header">
      <h1 class="page-title">System</h1>
      <span class="page-badge" v-if="fetchError">Demo data</span>
    </header>

    <!-- Loading -->
    <div class="loading-state" v-if="loading">
      <div class="spinner"></div>
      <span class="loading-text">Checking system…</span>
    </div>

    <template v-else>
      <!-- Ollama Section -->
      <section class="section">
        <div class="section-header">
          <h2 class="section-title">Ollama</h2>
          <span class="section-badge" :class="{ 'is-active': data.ollama?.loaded }">
            {{ data.ollama?.loaded ? 'Online' : 'Offline' }}
          </span>
        </div>

        <!-- VRAM gauge -->
        <div class="gauge-card">
          <div class="gauge-header">
            <span class="gauge-label">VRAM Usage</span>
            <span class="gauge-value">{{ data.ollama?.vram?.used ?? 0 }} / {{ data.ollama?.vram?.total ?? 0 }} GB</span>
          </div>
          <div class="gauge-track">
            <div
              class="gauge-fill"
              :style="{
                width: vramPercent + '%',
                background: gaugeColor(vramPercent),
              }"
            ></div>
          </div>
          <span class="gauge-percent" :style="{ color: gaugeColor(vramPercent) }">{{ vramPercent }}%</span>
        </div>

        <!-- Loaded models -->
        <div class="models-list" v-if="data.ollama?.models?.length">
          <div class="model-item" v-for="model in data.ollama.models" :key="model.name">
            <div class="model-info">
              <span class="model-name">{{ model.name }}</span>
              <span class="model-size">{{ model.size }}</span>
            </div>
            <div class="model-vram">
              <div class="mini-track">
                <div
                  class="mini-fill"
                  :style="{
                    width: model.vram + '%',
                    background: gaugeColor(model.vram),
                  }"
                ></div>
              </div>
              <span class="mini-label" :style="{ color: gaugeColor(model.vram) }">{{ model.vram }}% VRAM</span>
            </div>
          </div>
        </div>
        <div class="empty-state-inline" v-else>
          <span class="empty-text">No models loaded</span>
        </div>
      </section>

      <!-- System Resources Section -->
      <section class="section">
        <h2 class="section-title">System Resources</h2>

        <!-- CPU -->
        <div class="gauge-card">
          <div class="gauge-header">
            <span class="gauge-label">CPU</span>
            <span class="gauge-value">{{ cpuPercent }}% — {{ data.system?.cpu?.cores ?? 0 }} cores</span>
          </div>
          <div class="gauge-track">
            <div
              class="gauge-fill"
              :style="{
                width: cpuPercent + '%',
                background: gaugeColor(cpuPercent),
              }"
            ></div>
          </div>
        </div>

        <!-- RAM -->
        <div class="gauge-card">
          <div class="gauge-header">
            <span class="gauge-label">RAM</span>
            <span class="gauge-value">{{ data.system?.ram?.used ?? 0 }} / {{ data.system?.ram?.total ?? 0 }} GB</span>
          </div>
          <div class="gauge-track">
            <div
              class="gauge-fill"
              :style="{
                width: ramPercent + '%',
                background: gaugeColor(ramPercent),
              }"
            ></div>
          </div>
          <span class="gauge-percent" :style="{ color: gaugeColor(ramPercent) }">{{ ramPercent }}%</span>
        </div>

        <!-- Disk -->
        <div class="gauge-card">
          <div class="gauge-header">
            <span class="gauge-label">Disk</span>
            <span class="gauge-value">{{ data.system?.disk?.used ?? 0 }} / {{ data.system?.disk?.total ?? 0 }} GB</span>
          </div>
          <div class="gauge-track">
            <div
              class="gauge-fill"
              :style="{
                width: diskPercent + '%',
                background: gaugeColor(diskPercent),
              }"
            ></div>
          </div>
          <span class="gauge-percent" :style="{ color: gaugeColor(diskPercent) }">{{ diskPercent }}%</span>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.system-page {
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
  color: var(--color-warning);
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: var(--radius-full);
  padding: 2px var(--space-2);
}

/* Sections */
.section {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  margin-bottom: var(--space-6);
}

.section-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.section-title {
  font-family: var(--font-primary);
  font-size: var(--font-small);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-5) 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-header .section-title {
  margin-bottom: 0;
}

.section-badge {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  font-weight: 500;
  color: var(--color-text-muted);
  background: var(--color-surface);
  border-radius: var(--radius-full);
  padding: 2px var(--space-2);
}

.section-badge.is-active {
  color: var(--color-success);
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

/* Gauge cards */
.gauge-card {
  background: var(--color-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-5);
  margin-bottom: var(--space-3);
}

.gauge-card:last-child {
  margin-bottom: 0;
}

.gauge-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.gauge-label {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  font-weight: 600;
  color: var(--color-text-primary);
}

.gauge-value {
  font-family: var(--font-mono);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.gauge-track {
  height: 8px;
  background: var(--color-elevated);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.gauge-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width var(--transition-normal), background var(--transition-fast);
}

.gauge-percent {
  display: block;
  font-family: var(--font-mono);
  font-size: var(--font-micro);
  margin-top: var(--space-2);
  text-align: right;
}

/* Models list */
.models-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.model-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  background: var(--color-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  gap: var(--space-4);
}

.model-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 0;
}

.model-name {
  font-family: var(--font-mono);
  font-size: var(--font-caption);
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.model-size {
  font-family: var(--font-mono);
  font-size: var(--font-micro);
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.model-vram {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-shrink: 0;
}

.mini-track {
  width: 60px;
  height: 4px;
  background: var(--color-elevated);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.mini-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width var(--transition-normal);
}

.mini-label {
  font-family: var(--font-mono);
  font-size: var(--font-micro);
  white-space: nowrap;
}

.empty-state-inline {
  padding: var(--space-4) 0;
  text-align: center;
}

.empty-text {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

/* Loading */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-16) 0;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-standard);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-family: var(--font-primary);
  font-size: var(--font-small);
  color: var(--color-text-muted);
}

@media (max-width: 768px) {
  .system-page {
    padding: var(--space-4);
  }

  .model-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>