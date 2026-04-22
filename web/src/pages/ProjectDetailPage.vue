<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjects, useTasks } from '@/composables/useApi'

const route = useRoute()
const router = useRouter()
const projectId = computed(() => route.params.id)

const { projects, loading: projectsLoading, fetchAll } = useProjects()
const { tasks, loading: tasksLoading, fetchAll: fetchTasks } = useTasks()

const project = ref(null)

onMounted(async () => {
  await fetchAll()
  project.value = projects.value.find(p => p.id === Number(projectId.value) || p.id === projectId.value) || null
  await fetchTasks(projectId.value)
})

function goBack() {
  router.push('/projects')
}

function projectStatusColor(status) {
  const map = {
    active: 'var(--color-success)',
    paused: 'var(--color-warning)',
    completed: 'var(--color-accent)',
    archived: 'var(--color-text-muted)',
  }
  return map[status] || 'var(--color-text-muted)'
}

function taskStatusColor(status) {
  const map = {
    todo: 'var(--color-text-muted)',
    'in-progress': 'var(--color-info)',
    done: 'var(--color-success)',
    blocked: 'var(--color-error)',
    completed: 'var(--color-success)',
  }
  return map[status] || 'var(--color-text-muted)'
}

function progressWidth(p) {
  const v = p?.progress ?? p?.completion ?? 0
  return `${Math.min(100, Math.max(0, v))}%`
}

function taskCount() {
  return tasks.value?.length ?? project.value?.taskCount ?? 0
}

function completedCount() {
  return tasks.value?.filter(t => t.status === 'done' || t.status === 'completed').length ?? 0
}
</script>

<template>
  <div class="project-detail-page">
    <!-- Back button -->
    <button class="back-btn" @click="goBack">
      <span class="back-arrow">←</span>
      Back to Projects
    </button>

    <!-- Loading -->
    <div class="loading-state" v-if="projectsLoading">
      <div class="spinner"></div>
      <span class="loading-text">Loading project…</span>
    </div>

    <!-- Project found -->
    <template v-else-if="project">
      <header class="page-header">
        <div class="header-left">
          <h1 class="page-title">{{ project.name }}</h1>
          <span
            class="status-badge"
            :style="{ '--badge-color': projectStatusColor(project.status) }"
          >
            {{ project.status || 'unknown' }}
          </span>
        </div>
      </header>

      <!-- Stats row -->
      <section class="stats-row">
        <div class="stat-item">
          <span class="stat-value">{{ taskCount() }}</span>
          <span class="stat-label">Tasks</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ completedCount() }}</span>
          <span class="stat-label">Completed</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ progressWidth(project) }}</span>
          <span class="stat-label">Progress</span>
        </div>
      </section>

      <!-- Progress section -->
      <section class="progress-section">
        <div class="progress-header">
          <span class="progress-label">Overall Progress</span>
          <span class="progress-value">{{ progressWidth(project) }}</span>
        </div>
        <div class="progress-track">
          <div
            class="progress-fill"
            :style="{ width: progressWidth(project) }"
          ></div>
        </div>
      </section>

      <!-- Description -->
      <section class="description-section" v-if="project.description">
        <h2 class="section-title">Description</h2>
        <p class="description-text">{{ project.description }}</p>
      </section>

      <!-- Tasks list -->
      <section class="tasks-section">
        <div class="tasks-header">
          <h2 class="section-title">Tasks</h2>
          <span class="task-count-badge" v-if="tasks.length">{{ tasks.length }}</span>
        </div>

        <!-- Tasks loading -->
        <div class="loading-inline" v-if="tasksLoading">
          <div class="spinner-sm"></div>
          <span>Loading tasks…</span>
        </div>

        <!-- Tasks list -->
        <div class="tasks-list" v-else-if="tasks.length">
          <div class="task-item" v-for="task in tasks" :key="task.id">
            <span
              class="task-dot"
              :style="{ background: taskStatusColor(task.status) }"
            ></span>
            <div class="task-info">
              <span class="task-name">{{ task.name || task.title }}</span>
              <span
                class="task-status"
                :style="{ color: taskStatusColor(task.status) }"
              >{{ task.status }}</span>
            </div>
          </div>
        </div>

        <!-- Empty tasks -->
        <div class="empty-tasks" v-else>
          <span class="empty-text">No tasks in this project yet</span>
        </div>
      </section>
    </template>

    <!-- Not found -->
    <div class="not-found" v-else>
      <div class="not-found-icon">🔍</div>
      <h3 class="not-found-title">Project not found</h3>
      <p class="not-found-text">This project doesn't exist or has been removed.</p>
      <button class="btn-ghost" @click="goBack">Back to Projects</button>
    </div>
  </div>
</template>

<style scoped>
.project-detail-page {
  padding: var(--space-8);
  max-width: 1120px;
  margin: 0 auto;
}

/* Back button */
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  font-weight: 500;
  color: var(--color-text-muted);
  background: none;
  border: none;
  padding: 0;
  margin-bottom: var(--space-6);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.back-btn:hover {
  color: var(--color-text-secondary);
}

.back-arrow {
  font-size: 14px;
}

/* Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-6);
  padding-bottom: var(--space-5);
  border-bottom: 1px solid var(--border-standard);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.page-title {
  font-family: var(--font-primary);
  font-size: var(--font-h2);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
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
}

/* Stats row */
.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.stat-item {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.stat-value {
  font-family: var(--font-primary);
  font-size: var(--font-h3);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.stat-label {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

/* Progress */
.progress-section {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-5);
  margin-bottom: var(--space-6);
}

.progress-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.progress-label {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
}

.progress-value {
  font-family: var(--font-mono);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.progress-track {
  height: 6px;
  background: var(--color-surface);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--color-accent);
  border-radius: var(--radius-full);
  transition: width var(--transition-normal);
}

/* Description */
.description-section {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  margin-bottom: var(--space-6);
}

.section-title {
  font-family: var(--font-primary);
  font-size: var(--font-small);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-3) 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.description-text {
  font-family: var(--font-primary);
  font-size: var(--font-body);
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

/* Tasks section */
.tasks-section {
  background: var(--color-panel);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
}

.tasks-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.task-count-badge {
  font-family: var(--font-mono);
  font-size: var(--font-micro);
  color: var(--color-text-muted);
  background: var(--color-surface);
  border-radius: var(--radius-full);
  padding: 1px var(--space-2);
}

.tasks-list {
  display: flex;
  flex-direction: column;
}

.task-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) 0;
  border-bottom: 1px solid var(--border-subtle);
  transition: background var(--transition-fast);
}

.task-item:last-child {
  border-bottom: none;
}

.task-item:hover {
  background: rgba(255, 255, 255, 0.02);
  margin: 0 calc(var(--space-5) * -1);
  padding: var(--space-3) var(--space-5);
  border-radius: var(--radius-md);
}

.task-dot {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  flex-shrink: 0;
}

.task-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 0;
  flex: 1;
}

.task-name {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-status {
  font-family: var(--font-primary);
  font-size: var(--font-micro);
  text-transform: capitalize;
  flex-shrink: 0;
}

.loading-inline {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-4) 0;
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

.spinner-sm {
  width: 14px;
  height: 14px;
  border: 2px solid var(--border-standard);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.6s linear infinite;
}

.empty-tasks {
  padding: var(--space-8) 0;
  text-align: center;
}

.empty-text {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
}

/* Not found */
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-16) 0;
  text-align: center;
}

.not-found-icon {
  font-size: 48px;
}

.not-found-title {
  font-family: var(--font-primary);
  font-size: var(--font-h3);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.not-found-text {
  font-family: var(--font-primary);
  font-size: var(--font-caption);
  color: var(--color-text-muted);
  margin: 0;
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
  margin-top: var(--space-2);
  transition: all var(--transition-fast);
}

.btn-ghost:hover {
  background: var(--color-elevated);
  border-color: var(--border-strong);
}

/* Loading state */
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
  .project-detail-page {
    padding: var(--space-4);
  }

  .stats-row {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }
}
</style>