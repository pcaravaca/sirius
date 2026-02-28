<template>
  <div class="app">
    <header class="header">
      <div class="logo">
        <span class="logo-icon">⭐</span>
        <h1>Sirius</h1>
      </div>
      <div class="header-info">
        <span class="weather" v-if="dashboard.weather">
          {{ dashboard.weather.icon }} {{ dashboard.weather.temp }}°C
        </span>
        <span class="exchange" v-if="dashboard.exchange">
          💵 ₡{{ dashboard.exchange.usdToCrc.toFixed(2) }} / ₡{{ dashboard.exchange.sellRate.toFixed(2) }}
        </span>
      </div>
    </header>

    <main class="main">
      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">📁</div>
          <div class="stat-content">
            <span class="stat-value">{{ dashboard.projects?.length || 0 }}</span>
            <span class="stat-label">Proyectos</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <span class="stat-value">{{ completedTasks }}</span>
            <span class="stat-label">Completadas</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⏳</div>
          <div class="stat-content">
            <span class="stat-value">{{ pendingTasks }}</span>
            <span class="stat-label">Pendientes</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">🔥</div>
          <div class="stat-content">
            <span class="stat-value">{{ activeProjects }}</span>
            <span class="stat-label">Activos</span>
          </div>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="charts-row">
        <div class="chart-card">
          <h3>Progreso de Proyectos</h3>
          <div class="chart-container">
            <canvas ref="progressChart"></canvas>
          </div>
        </div>
        <div class="chart-card">
          <h3>Estado de Tareas</h3>
          <div class="chart-container">
            <canvas ref="taskChart"></canvas>
          </div>
        </div>
      </div>

      <!-- Tareas por Proyecto -->
      <div class="projects-tasks-section">
        <h2>📋 Tareas por Proyecto</h2>
        
        <div class="project-tabs">
          <button 
            v-for="project in dashboard.projects" 
            :key="project.id"
            class="project-tab"
            :class="{ active: activeProjectTab === project.id }"
            @click="activeProjectTab = project.id"
          >
            {{ project.name }}
            <span class="task-count">
              {{ getProjectTaskCount(project.id) }}
            </span>
          </button>
        </div>

        <div class="project-tasks-content" v-if="activeProject">
          <div class="tasks-header">
            <h3>{{ activeProject.name }}</h3>
            <span class="project-status" :class="activeProject.status">
              {{ activeProject.status }}
            </span>
          </div>
          
          <div class="project-progress">
            <div class="progress-bar large">
              <div class="progress-fill" :style="{ width: activeProject.progress + '%' }"></div>
            </div>
            <span class="progress-text">{{ activeProject.progress }}%</span>
          </div>

          <div class="tasks-list">
            <div 
              v-for="task in getProjectTasks(activeProject.id)" 
              :key="task.id" 
              class="task-row"
              :class="task.status"
            >
              <input 
                type="checkbox" 
                :checked="task.status === 'completed'"
                @change="toggleTask(task.id)"
              />
              <div class="task-info">
                <span class="task-title" :class="{ completed: task.status === 'completed' }">
                  {{ task.title }}
                </span>
                <span class="task-desc" v-if="task.description">{{ task.description }}</span>
              </div>
              <span class="task-due" v-if="task.dueDate">{{ task.dueDate }}</span>
              <span class="priority-badge" :class="task.priority">{{ task.priority }}</span>
              <span class="task-status-badge" :class="task.status">{{ task.status }}</span>
            </div>
            
            <div class="no-tasks" v-if="getProjectTasks(activeProject.id).length === 0">
              No hay tareas para este proyecto
            </div>
          </div>
        </div>
      </div>

      <!-- Motivational -->
      <div class="motivational-card" v-if="dashboard.motivational">
        <span>💬</span>
        <p>{{ dashboard.motivational }}</p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import axios from 'axios'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

const API_URL = '/api'

const dashboard = ref({})
const progressChart = ref(null)
const taskChart = ref(null)
const activeProjectTab = ref(null)

const completedTasks = computed(() => 
  dashboard.value.tasks?.filter(t => t.status === 'completed').length || 0
)

const pendingTasks = computed(() => 
  dashboard.value.tasks?.filter(t => t.status !== 'completed').length || 0
)

const activeProjects = computed(() => 
  dashboard.value.projects?.filter(p => p.status === 'active').length || 0
)

const activeProject = computed(() => {
  if (!activeProjectTab.value) return null
  return dashboard.value.projects?.find(p => p.id === activeProjectTab.value)
})

const getProjectName = (id) => {
  const p = dashboard.value.projects?.find(p => p.id === id)
  return p?.name || 'Sin proyecto'
}

const getProjectTasks = (projectId) => {
  return dashboard.value.tasks?.filter(t => t.projectId === projectId) || []
}

const getProjectTaskCount = (projectId) => {
  const tasks = getProjectTasks(projectId)
  const completed = tasks.filter(t => t.status === 'completed').length
  return `${completed}/${tasks.length}`
}

const fetchDashboard = async () => {
  try {
    const { data } = await axios.get(`${API_URL}/dashboard`)
    dashboard.value = data
    
    // Select first project by default
    if (!activeProjectTab.value && data.projects?.length > 0) {
      activeProjectTab.value = data.projects[0].id
    }
    
    await nextTick()
    renderCharts()
  } catch (e) {
    console.error('Error fetching dashboard:', e)
  }
}

const toggleTask = async (id) => {
  try {
    await axios.patch(`${API_URL}/tasks/${id}/toggle`)
    fetchDashboard()
  } catch (e) {
    console.error('Error toggling task:', e)
  }
}

let progressChartInstance = null
let taskChartInstance = null

const renderCharts = () => {
  if (!dashboard.value.projects) return

  if (progressChart.value) {
    if (progressChartInstance) progressChartInstance.destroy()
    
    const projects = dashboard.value.projects || []
    progressChartInstance = new Chart(progressChart.value, {
      type: 'bar',
      data: {
        labels: projects.map(p => p.name),
        datasets: [{
          label: 'Progreso %',
          data: projects.map(p => p.progress),
          backgroundColor: ['#8b5cf6', '#06b6d4', '#10b981', '#f59e0b', '#ef4444'],
          borderRadius: 8
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          y: { beginAtZero: true, max: 100, grid: { color: '#1f2937' } },
          x: { grid: { display: false } }
        }
      }
    })
  }

  if (taskChart.value) {
    if (taskChartInstance) taskChartInstance.destroy()
    
    const tasks = dashboard.value.tasks || []
    const completed = tasks.filter(t => t.status === 'completed').length
    const pending = tasks.filter(t => t.status !== 'completed').length

    taskChartInstance = new Chart(taskChart.value, {
      type: 'doughnut',
      data: {
        labels: ['Completadas', 'Pendientes'],
        datasets: [{
          data: [completed, pending],
          backgroundColor: ['#10b981', '#f59e0b'],
          borderWidth: 0
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { position: 'bottom', labels: { color: '#9ca3af' } } }
      }
    })
  }
}

onMounted(() => {
  fetchDashboard()
  
  const ws = new WebSocket(`ws://${window.location.host}/ws`)
  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data)
    if (['project_created', 'project_updated', 'project_deleted', 
         'task_created', 'task_updated', 'task_toggled', 'task_deleted'].includes(msg.type)) {
      fetchDashboard()
    }
  }
})
</script>

<style>
.app {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 100%);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: rgba(255,255,255,0.05);
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.logo-icon {
  font-size: 1.5rem;
}

.logo h1 {
  font-size: 1.5rem;
  font-weight: 700;
  background: linear-gradient(90deg, #8b5cf6, #06b6d4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-info {
  display: flex;
  gap: 1.5rem;
}

.weather, .exchange {
  padding: 0.5rem 1rem;
  background: rgba(255,255,255,0.1);
  border-radius: 8px;
  font-size: 0.9rem;
}

.main {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: rgba(255,255,255,0.05);
  border-radius: 12px;
  border: 1px solid rgba(255,255,255,0.1);
}

.stat-icon {
  font-size: 2rem;
}

.stat-value {
  display: block;
  font-size: 2rem;
  font-weight: 700;
}

.stat-label {
  color: #9ca3af;
  font-size: 0.85rem;
}

.charts-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.chart-card {
  background: rgba(255,255,255,0.05);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid rgba(255,255,255,0.1);
}

.chart-card h3 {
  margin-bottom: 1rem;
  color: #e4e4e7;
}

.chart-container {
  height: 250px;
}

/* Projects Tasks Section */
.projects-tasks-section {
  background: rgba(255,255,255,0.05);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid rgba(255,255,255,0.1);
  margin-bottom: 2rem;
}

.projects-tasks-section h2 {
  margin-bottom: 1rem;
  color: #e4e4e7;
}

.project-tabs {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 1.5rem;
}

.project-tab {
  padding: 0.75rem 1rem;
  background: rgba(255,255,255,0.1);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.project-tab:hover {
  background: rgba(255,255,255,0.15);
}

.project-tab.active {
  background: #8b5cf6;
  color: white;
  border-color: #8b5cf6;
}

.task-count {
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  background: rgba(0,0,0,0.2);
  border-radius: 10px;
}

.project-tab.active .task-count {
  background: rgba(0,0,0,0.3);
}

.project-tasks-content {
  background: rgba(0,0,0,0.2);
  border-radius: 8px;
  padding: 1.5rem;
}

.tasks-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.tasks-header h3 {
  color: #e4e4e7;
  margin: 0;
}

.project-status {
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.75rem;
  text-transform: uppercase;
}

.project-status.active { background: #10b981; color: white; }
.project-status.paused { background: #f59e0b; color: white; }
.project-status.completed { background: #6b7280; color: white; }

.project-progress {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.progress-bar.large {
  flex: 1;
  height: 12px;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: rgba(255,255,255,0.1);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #8b5cf6, #06b6d4);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 0.9rem;
  color: #9ca3af;
  min-width: 40px;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.task-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255,255,255,0.05);
  border-radius: 8px;
  border-left: 3px solid transparent;
}

.task-row.completed {
  border-left-color: #10b981;
}

.task-row.pending {
  border-left-color: #f59e0b;
}

.task-row input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.task-info {
  flex: 1;
}

.task-title {
  display: block;
  font-weight: 500;
  color: #e4e4e7;
}

.task-title.completed {
  text-decoration: line-through;
  color: #6b7280;
}

.task-desc {
  display: block;
  font-size: 0.85rem;
  color: #9ca3af;
}

.task-due {
  font-size: 0.85rem;
  color: #9ca3af;
}

.priority-badge, .task-status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.7rem;
  text-transform: uppercase;
}

.priority-badge.high { background: #ef4444; color: white; }
.priority-badge.medium { background: #f59e0b; color: white; }
.priority-badge.low { background: #6b7280; color: white; }

.task-status-badge.completed { background: #10b981; color: white; }
.task-status-badge.pending { background: #f59e0b; color: white; }
.task-status-badge.in_progress { background: #3b82f6; color: white; }

.no-tasks {
  text-align: center;
  padding: 2rem;
  color: #6b7280;
}

.motivational-card {
  background: linear-gradient(90deg, rgba(139, 92, 246, 0.2), rgba(6, 182, 212, 0.2));
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.motivational-card p {
  font-size: 1.1rem;
  font-style: italic;
}
</style>
