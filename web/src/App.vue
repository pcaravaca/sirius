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
            <span class="stat-label">Tareas Completadas</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⏳</div>
          <div class="stat-content">
            <span class="stat-value">{{ pendingTasks }}</span>
            <span class="stat-label">Tareas Pendientes</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">🔥</div>
          <div class="stat-content">
            <span class="stat-value">{{ activeProjects }}</span>
            <span class="stat-label">Proyectos Activos</span>
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

      <!-- Projects & Tasks -->
      <div class="content-row">
        <div class="card projects-card">
          <div class="card-header">
            <h2>📁 Proyectos</h2>
            <button class="btn-add" @click="showProjectModal = true">+</button>
          </div>
          <div class="list">
            <div v-for="project in dashboard.projects" :key="project.id" class="item">
              <div class="item-info">
                <span class="item-name">{{ project.name }}</span>
                <span class="item-desc">{{ project.description }}</span>
              </div>
              <div class="item-status" :class="project.status">
                {{ project.status }}
              </div>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: project.progress + '%' }"></div>
              </div>
            </div>
          </div>
        </div>

        <div class="card tasks-card">
          <div class="card-header">
            <h2>✅ Tareas</h2>
            <button class="btn-add" @click="showTaskModal = true">+</button>
          </div>
          <div class="list">
            <div v-for="task in dashboard.tasks" :key="task.id" class="item task-item">
              <input 
                type="checkbox" 
                :checked="task.status === 'completed'"
                @change="toggleTask(task.id)"
              />
              <div class="item-info">
                <span class="item-title" :class="{ completed: task.status === 'completed' }">
                  {{ task.title }}
                </span>
                <span class="item-meta">
                  {{ getProjectName(task.projectId) }} • {{ task.dueDate }}
                </span>
              </div>
              <span class="priority" :class="task.priority">{{ task.priority }}</span>
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

    <!-- Modals -->
    <div class="modal" v-if="showProjectModal" @click.self="showProjectModal = false">
      <div class="modal-content">
        <h3>Nuevo Proyecto</h3>
        <input v-model="newProject.name" placeholder="Nombre" />
        <input v-model="newProject.description" placeholder="Descripción" />
        <select v-model="newProject.status">
          <option value="active">Activo</option>
          <option value="paused">Pausado</option>
          <option value="completed">Completado</option>
        </select>
        <select v-model="newProject.priority">
          <option value="high">Alta</option>
          <option value="medium">Media</option>
          <option value="low">Baja</option>
        </select>
        <div class="modal-actions">
          <button @click="createProject">Crear</button>
          <button class="cancel" @click="showProjectModal = false">Cancelar</button>
        </div>
      </div>
    </div>

    <div class="modal" v-if="showTaskModal" @click.self="showTaskModal = false">
      <div class="modal-content">
        <h3>Nueva Tarea</h3>
        <input v-model="newTask.title" placeholder="Título" />
        <input v-model="newTask.description" placeholder="Descripción" />
        <select v-model="newTask.projectId">
          <option value="">Seleccionar Proyecto</option>
          <option v-for="p in dashboard.projects" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>
        <select v-model="newTask.priority">
          <option value="high">Alta</option>
          <option value="medium">Media</option>
          <option value="low">Baja</option>
        </select>
        <input v-model="newTask.dueDate" type="date" />
        <div class="modal-actions">
          <button @click="createTask">Crear</button>
          <button class="cancel" @click="showTaskModal = false">Cancelar</button>
        </div>
      </div>
    </div>
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
const showProjectModal = ref(false)
const showTaskModal = ref(false)

const newProject = ref({ name: '', description: '', status: 'active', priority: 'medium' })
const newTask = ref({ title: '', description: '', projectId: '', priority: 'medium', dueDate: '' })

const completedTasks = computed(() => 
  dashboard.value.tasks?.filter(t => t.status === 'completed').length || 0
)

const pendingTasks = computed(() => 
  dashboard.value.tasks?.filter(t => t.status !== 'completed').length || 0
)

const activeProjects = computed(() => 
  dashboard.value.projects?.filter(p => p.status === 'active').length || 0
)

const getProjectName = (id) => {
  const p = dashboard.value.projects?.find(p => p.id === id)
  return p?.name || 'Sin proyecto'
}

const fetchDashboard = async () => {
  try {
    const { data } = await axios.get(`${API_URL}/dashboard`)
    dashboard.value = data
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

const createProject = async () => {
  try {
    await axios.post(`${API_URL}/projects`, { ...newProject.value, progress: 0 })
    showProjectModal.value = false
    newProject.value = { name: '', description: '', status: 'active', priority: 'medium' }
    fetchDashboard()
  } catch (e) {
    console.error('Error creating project:', e)
  }
}

const createTask = async () => {
  try {
    await axios.post(`${API_URL}/tasks`, { ...newTask.value, status: 'pending' })
    showTaskModal.value = false
    newTask.value = { title: '', description: '', projectId: '', priority: 'medium', dueDate: '' }
    fetchDashboard()
  } catch (e) {
    console.error('Error creating task:', e)
  }
}

let progressChartInstance = null
let taskChartInstance = null

const renderCharts = () => {
  if (!dashboard.value.projects) return

  // Progress Chart (Bar)
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

  // Task Chart (Doughnut)
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
  
  // WebSocket for real-time updates
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

.content-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.card {
  background: rgba(255,255,255,0.05);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid rgba(255,255,255,0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.card-header h2 {
  font-size: 1.25rem;
}

.btn-add {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: #8b5cf6;
  color: white;
  font-size: 1.25rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-add:hover {
  background: #7c3aed;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(0,0,0,0.2);
  border-radius: 8px;
}

.item-info {
  flex: 1;
}

.item-name {
  display: block;
  font-weight: 500;
}

.item-desc {
  display: block;
  color: #9ca3af;
  font-size: 0.85rem;
}

.item-title.completed {
  text-decoration: line-through;
  color: #9ca3af;
}

.item-meta {
  color: #6b7280;
  font-size: 0.8rem;
}

.item-status {
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.75rem;
  text-transform: uppercase;
}

.item-status.active { background: #10b981; color: white; }
.item-status.paused { background: #f59e0b; color: white; }
.item-status.completed { background: #6b7280; color: white; }

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

.priority {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.7rem;
  text-transform: uppercase;
}

.priority.high { background: #ef4444; color: white; }
.priority.medium { background: #f59e0b; color: white; }
.priority.low { background: #6b7280; color: white; }

.task-item input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
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

/* Modal */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal-content {
  background: #1a1a2e;
  padding: 2rem;
  border-radius: 12px;
  width: 400px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.modal-content h3 {
  margin-bottom: 0.5rem;
}

.modal-content input,
.modal-content select {
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.2);
  background: rgba(255,255,255,0.1);
  color: white;
  font-size: 1rem;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.modal-actions button {
  flex: 1;
  padding: 0.75rem;
  border-radius: 8px;
  border: none;
  background: #8b5cf6;
  color: white;
  font-size: 1rem;
  cursor: pointer;
}

.modal-actions button.cancel {
  background: rgba(255,255,255,0.2);
}

.modal-actions button:hover {
  opacity: 0.9;
}
</style>
