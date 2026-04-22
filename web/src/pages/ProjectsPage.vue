<script setup>
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useProjects } from '@/composables/useApi'

const router = useRouter()
const { projects, loading, fetchAll, create } = useProjects()

const showCreateModal = ref(false)
const newProjectName = ref('')
const creating = ref(false)

onMounted(() => {
  fetchAll()
})

function statusColor(status) {
  const map = {
    active: 'var(--color-success)',
    paused: 'var(--color-warning)',
    completed: 'var(--color-accent)',
    archived: 'var(--color-text-muted)',
  }
  return map[status] || 'var(--color-text-muted)'
}

function progressWidth(project) {
  const p = project.progress ?? project.completion ?? 0
  return `${Math.min(100, Math.max(0, p))}%`
}

function taskCount(project) {
  return project.taskCount ?? project.tasks ?? 0
}

function openProject(id) {
  router.push(`/projects/${id}`)
}

async function handleCreate() {
  if (!newProjectName.value.trim()) return
  creating.value = true
  try {
    await create({ name: newProjectName.value.trim() })
    newProjectName.value = ''
    showCreateModal.value = false
    await fetchAll()
  } catch (e) {
    console.error('Failed to create project', e)
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <div class="projects-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">Projects</h1>
        <span class="page-count" v-if="projects.length">{{ projects.length }} total</span>
      </div>
      <button class="btn-primary" @click="showCreateModal = true">
        <span class="btn-icon">+</span>
        New Project
      </button>
    </header>

    <!-- Create Modal -->
    <div class="modal-overlay" v-if="showCreateModal" @click.self="showCreateModal = false">
      <div class="modal-panel">
        <h3 class="modal-title">Create Project</h3>
        <input
          class="modal-input"
          v-model="newProjectName"
          placeholder="Project name"
          @keyup.enter="handleCreate"
          autofocus
        />
        <div class="modal-actions">
          <button class="btn-ghost" @click="showCreateModal = false">Cancel</button>
          <button class="btn-primary" @click="handleCreate" :disabled="creating">
            {{ creating ? 'Creating…' : 'Create' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div class="loading-state" v-if="loading">
      <div class="spinner"></div>
      <span class="loading-text">Loading projects…</span>
    </div>

    <!-- Grid -->
    <div class="projects-grid" v-else-if="projects.length">
      <div
        class="project-card"
        v-for="project in projects"
        :key="project.id"
        @click="openProject(project.id)"
      >
        <div class="card-top">
          <h3 class="project-name">{{ project.name }}</h3>
          <span
            class="status-badge"
            :style="{ '--badge-color': statusColor(project.status) }"
          >
            {{ project.status || 'unknown' }}
          </span>
        </div>
        <p class="project-desc" v-if="project.description">{{ project.description }}</p>
        <div class="progress-track">
          <div
            class="progress-fill"
            :style="{ width: progressWidth(project), '--fill-color': statusColor(project.status) }"
          ></div>
        </div>
        <div class="card-footer">
          <span class="progress-label">{{ progressWidth(project) }} complete</span>
          <span class="task-count">{{ taskCount(project) }} tasks</span>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div class="empty-state" v-else>
      <div class="empty-icon">📁</div>
      <h3 class="empty-title">No projects yet</h3>
      <p class="empty-text">Create your first project to get started.</p>
      <button class="btn-primary" @click="showCreateModal = true">
        <span class="btn-icon">+</span>
        New Project
      </button>
    </div>
  </div>
</template>

<style scoped>
.projects-page {
  padding: var(--space-8);
  max-width: 1120px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-8);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
}

.page-title {
  font-family: var(--font-primary);
  font-size: var(--font-h2);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.page-count {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  font-weight: 500;
  color: var(--color-text-primary);
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius-lg);
  padding: var(--space-2) var(--space-4);
  cursor: pointer;
  transition: background var(--transition-fast), opacity var(--transition-fast);
}

.btn-primary:hover {
  background: var(--color-accent-hover);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  font-size: 16px;
  font-weight: 600;
  line-height: 1;
}

.btn-ghost {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  font-weight: 500;
  color: var(--color-text-secondary);
  background: none;
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-2) var(--space-4);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-ghost:hover {
  background: var(--color-elevated);
  border-color: var(--border-strong);
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal-panel {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  width: 400px;
  max-width: 90vw;
  box-shadow: var(--shadow-lg);
}

.modal-title {
  font-family: var(--font-primary);
  font-size: var(--font-h3);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 var(--space-4) 0;
}

.modal-input {
  width: 100%;
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-2) var(--space-3);
  font-family: var(--font-primary);
  font-size: var(--font-body);
  color: var(--color-text-primary);
  outline: none;
  transition: border-color var(--transition-fast);
  box-sizing: border-box;
}

.modal-input:focus {
  border-color: var(--color-accent);
}

.modal-input::placeholder {
  color: var(--color-text-muted);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-5);
}

/* Grid */
.projects-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
}

.project-card {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  cursor: pointer;
  transition: border-color var(--transition-fast), background var(--transition-fast);
}

.project-card:hover {
  border-color: var(--border-strong);
  background: var(--color-surface);
}

.card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
}

.project-name {
  font-family: var(--font-primary);
  font-size: var(--font-body);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-badge {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  font-weight: 500;
  color: var(--badge-color);
  background: color-mix(in srgb, var(--badge-color) 12%, transparent);
  border-radius: var(--radius-full);
  padding: 2px var(--space-2);
  text-transform: capitalize;
  white-space: nowrap;
  flex-shrink: 0;
}

.project-desc {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
  margin: 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.progress-track {
  height: 4px;
  background: var(--color-surface);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--fill-color, var(--color-accent));
  border-radius: var(--radius-full);
  transition: width var(--transition-normal);
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.progress-label,
.task-count {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  color: var(--color-text-muted);
}

/* Loading */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
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

/* Empty */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-16) 0;
}

.empty-icon {
  font-size: 40px;
}

.empty-title {
  font-family: var(--font-primary);
  font-size: var(--font-h3);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.empty-text {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
  margin: 0;
}

@media (max-width: 1024px) {
  .projects-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .projects-page {
    padding: var(--space-4);
  }

  .projects-grid {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }

  .modal-panel {
    width: 95vw;
  }
}
</style>