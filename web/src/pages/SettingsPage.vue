<template>
  <div class="settings-page">
    <div class="page-header">
      <h2>⚙️ Settings</h2>
    </div>

    <div class="settings-sections">
      <section class="settings-section">
        <h3>Services</h3>
        <div class="setting-row">
          <label>OpenClaw Gateway</label>
          <input v-model="settings.openclawUrl" type="url" placeholder="http://localhost:18789" />
        </div>
        <div class="setting-row">
          <label>Ollama Server</label>
          <input v-model="settings.ollamaUrl" type="url" placeholder="http://localhost:11434" />
        </div>
        <div class="setting-row">
          <label>Poll Interval</label>
          <select v-model="settings.pollInterval">
            <option value="10">10s</option>
            <option value="30">30s</option>
            <option value="60">60s</option>
          </select>
        </div>
      </section>

      <section class="settings-section">
        <h3>Agent Types</h3>
        <div class="setting-row">
          <label>OpenClaw</label>
          <span class="status-dot on"></span>
          <span class="status-label">Enabled</span>
        </div>
        <div class="setting-row">
          <label>Claude Code</label>
          <input v-model="settings.claudePath" type="text" placeholder="/usr/local/bin/claude" />
        </div>
        <div class="setting-row">
          <label>OpenAI Codex</label>
          <input v-model="settings.codexPath" type="text" placeholder="/usr/local/bin/codex" />
        </div>
        <div class="setting-row">
          <label>Hermes Agent</label>
          <span class="status-dot on"></span>
          <span class="status-label">Built-in</span>
        </div>
      </section>

      <section class="settings-section">
        <h3>Display</h3>
        <div class="setting-row">
          <label>Theme</label>
          <div class="toggle-group">
            <button class="toggle-btn active">🌑 Dark</button>
            <button class="toggle-btn">☀️ Light</button>
          </div>
        </div>
      </section>

      <div class="actions">
        <button class="btn-primary" @click="save">Save Settings</button>
        <button class="btn-ghost" @click="reset">Reset</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const settings = reactive({
  openclawUrl: localStorage.getItem('sirius_openclaw') || 'http://localhost:18789',
  ollamaUrl: localStorage.getItem('sirius_ollama') || 'http://localhost:11434',
  pollInterval: localStorage.getItem('sirius_poll') || '30',
  claudePath: localStorage.getItem('sirius_claude') || '',
  codexPath: localStorage.getItem('sirius_codex') || '',
})

function save() {
  localStorage.setItem('sirius_openclaw', settings.openclawUrl)
  localStorage.setItem('sirius_ollama', settings.ollamaUrl)
  localStorage.setItem('sirius_poll', settings.pollInterval)
  localStorage.setItem('sirius_claude', settings.claudePath)
  localStorage.setItem('sirius_codex', settings.codexPath)
}

function reset() {
  settings.openclawUrl = 'http://localhost:18789'
  settings.ollamaUrl = 'http://localhost:11434'
  settings.pollInterval = '30'
  settings.claudePath = ''
  settings.codexPath = ''
}
</script>

<style scoped>
.settings-page {
  max-width: 700px;
}

.page-header {
  margin-bottom: var(--space-6);
}

.page-header h2 {
  font-size: var(--font-h3);
  font-weight: 590;
}

.settings-sections {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.settings-section {
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
}

.settings-section h3 {
  font-size: var(--font-small);
  font-weight: 510;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--border-subtle);
}

.setting-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) 0;
}

.setting-row + .setting-row {
  border-top: 1px solid var(--border-subtle);
}

.setting-row label {
  flex: 0 0 160px;
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.setting-row input,
.setting-row select {
  flex: 1;
  background: var(--color-canvas);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
  color: var(--color-text-primary);
  font-size: var(--font-caption);
  font-family: var(--font-mono);
  outline: none;
}

.setting-row input:focus,
.setting-row select:focus {
  border-color: var(--color-accent);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-dot.on { background: var(--color-success); }

.status-label {
  font-size: var(--font-caption);
  color: var(--color-text-subtle);
}

.toggle-group {
  display: flex;
  background: var(--color-canvas);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.toggle-btn {
  padding: var(--space-2) var(--space-4);
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  font-size: var(--font-caption);
  cursor: pointer;
}

.toggle-btn.active {
  background: var(--color-accent);
  color: white;
}

.actions {
  display: flex;
  gap: var(--space-3);
  padding-top: var(--space-4);
}

.btn-primary {
  background: var(--color-accent);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-5);
  font-size: var(--font-caption);
  font-weight: 510;
  cursor: pointer;
}

.btn-primary:hover { background: var(--color-accent-bright); }

.btn-ghost {
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-5);
  color: var(--color-text-muted);
  font-size: var(--font-caption);
  cursor: pointer;
}

.btn-ghost:hover { background: var(--color-elevated); }

@media (max-width: 768px) {
  .setting-row {
    flex-direction: column;
    align-items: flex-start;
  }
  .setting-row label {
    flex: none;
    margin-bottom: var(--space-1);
  }
  .setting-row input,
  .setting-row select {
    width: 100%;
  }
}
</style>