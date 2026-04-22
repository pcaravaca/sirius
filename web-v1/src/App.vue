<template>
  <div class="app">
    <header class="header">
      <div class="logo">
        <span class="logo-icon">⭐</span>
        <h1>Sirius</h1>
      </div>
      <div class="header-info">
        <div class="weather-compact" v-if="dashboard.weather">
          <span>{{ dashboard.weather.icon }}</span>
          <span class="temp">{{ dashboard.weather.temp }}°C</span>
          <span class="condition">{{ dashboard.weather.condition }}</span>
        </div>
        <span class="exchange" v-if="dashboard.exchange">
          💵 ₡{{ dashboard.exchange.usdToCrc.toFixed(2) }}
        </span>
      </div>
    </header>

    <main class="main">
      <!-- Navigation Tabs -->
      <div class="nav-tabs">
        <button class="nav-tab" :class="{ active: activeTab === 'dashboard' }" @click="activeTab = 'dashboard'">
          📊 Dashboard
        </button>
        <button class="nav-tab" :class="{ active: activeTab === 'projects' }" @click="activeTab = 'projects'">
          📁 Proyectos
        </button>
        <button class="nav-tab" :class="{ active: activeTab === 'kanban' }" @click="activeTab = 'kanban'">
          📋 Kanban
        </button>
        <button class="nav-tab" :class="{ active: activeTab === 'agents' }" @click="activeTab = 'agents'; fetchAgentsStatus()">
          🤖 Agentes
        </button>
        <button class="nav-tab" :class="{ active: activeTab === 'exchange' }" @click="activeTab = 'exchange'; fetchExchangeData()">
          💱 Tipo de Cambio
        </button>
        <button class="nav-tab" :class="{ active: activeTab === 'news' }" @click="activeTab = 'news'; fetchNews()">
          📰 Noticias
        </button>
      </div>

      <!-- DASHBOARD TAB -->
      <div v-if="activeTab === 'dashboard'">
        <!-- Weather Card -->
        <div class="weather-section" v-if="dashboard.weather">
          <div class="weather-main">
            <div class="weather-icon">{{ dashboard.weather.icon }}</div>
            <div class="weather-temp">
              <span class="temp-value">{{ dashboard.weather.temp }}°</span>
              <span class="temp-feels">Sensación {{ dashboard.weather.feelsLike }}°</span>
            </div>
            <div class="weather-location">
              <span class="location-name">{{ dashboard.weather.location }}</span>
              <span class="condition-text">{{ dashboard.weather.condition }}</span>
            </div>
          </div>
          
          <div class="weather-details">
            <div class="weather-detail">
              <span class="detail-icon">💧</span>
              <span class="detail-label">Humedad</span>
              <span class="detail-value">{{ dashboard.weather.humidity }}%</span>
            </div>
            <div class="weather-detail">
              <span class="detail-icon">💨</span>
              <span class="detail-label">Viento</span>
              <span class="detail-value">{{ dashboard.weather.windSpeed }} km/h {{ dashboard.weather.windDir }}</span>
            </div>
            <div class="weather-detail">
              <span class="detail-icon">☁️</span>
              <span class="detail-label">Nubosidad</span>
              <span class="detail-value">{{ dashboard.weather.cloudCover }}%</span>
            </div>
            <div class="weather-detail">
              <span class="detail-icon">📊</span>
              <span class="detail-label">Presión</span>
              <span class="detail-value">{{ dashboard.weather.pressure }} hPa</span>
            </div>
          </div>

          <div class="weather-forecast" v-if="dashboard.weather.forecast && dashboard.weather.forecast.length > 0">
            <h4>Pronóstico 5 días</h4>
            <div class="forecast-grid">
              <div class="forecast-day" v-for="day in dashboard.weather.forecast" :key="day.date">
                <span class="forecast-date">{{ formatDate(day.date) }}</span>
                <span class="forecast-icon">{{ day.icon }}</span>
                <span class="forecast-temp">{{ day.tempMax }}° / {{ day.tempMin }}°</span>
                <span class="forecast-precip" v-if="day.precipitation > 0">🌧️ {{ day.precipitation }}mm</span>
              </div>
            </div>
          </div>
        </div>

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
              <span class="stat-value">{{ inProgressTasks }}</span>
              <span class="stat-label">En Progreso</span>
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
              <span class="task-count">{{ getProjectTaskCount(project.id) }}</span>
            </button>
          </div>

          <div class="project-tasks-content" v-if="activeProject">
            <div class="tasks-header">
              <h3>{{ activeProject.name }}</h3>
              <span class="project-status" :class="activeProject.status">{{ activeProject.status }}</span>
              
            </div>
            
            <div class="project-progress">
              <div class="progress-bar large">
                <div class="progress-fill" :style="{ width: activeProject.progress + '%' }"></div>
              </div>
              <span class="progress-text">{{ activeProject.progress }}%</span>
            </div>

            <div class="tasks-list">
              <div v-for="task in getProjectTasks(activeProject.id)" :key="task.id" class="task-row" :class="task.status">
                <div class="task-status-indicator">{{ getStatusIcon(task.status) }}</div>
                
                <div class="task-info">
                  <span class="task-title" :class="{ completed: task.status === 'completed' }">{{ task.title }}</span>
                  <span class="task-desc" v-if="task.description">{{ task.description }}</span>
                  <span class="task-due" v-if="task.dueDate">📅 {{ task.dueDate }}</span>
                </div>
                
                <div class="task-actions">
                  <select class="status-select" :class="task.status" :value="task.status" @change="setTaskStatus(task.id, $event.target.value)">
                    <option value="pending">⏳ Pendiente</option>
                    <option value="in_progress">🔥 En progreso</option>
                    <option value="paused">⏸️ Pausado</option>
                    <option value="completed">✅ Completado</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- PROJECTS TAB -->
      <div v-if="activeTab === 'projects'" class="projects-section">
        <div class="section-header">
          <h2>📁 Gestión de Proyectos</h2>
          <button @click="openProjectModal()" class="add-btn">➕ Nuevo Proyecto</button>
        </div>

        <div class="projects-grid">
          <div v-for="project in dashboard.projects" :key="project.id" class="project-card">
            <div class="project-header">
              <h3>{{ project.name }}</h3>
              <div class="project-actions">
                <button @click="openProjectModal(project)" class="icon-btn">✏️</button>
                <button @click="deleteProject(project.id)" class="icon-btn danger">🗑️</button>
              </div>
            </div>
            <p class="project-desc" v-if="project.description">{{ project.description }}</p>
            <div class="project-meta">
              <span class="status-badge" :class="project.status">{{ project.status }}</span>
              <span class="progress-label">{{ project.progress }}%</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: project.progress + '%' }"></div>
            </div>
            <div class="project-footer">
              <span>{{ getProjectTaskCount(project.id) }} tareas</span>
              <span>{{ getProjectCompletedCount(project.id) }}/{{ getProjectTasks(project.id).length }} completadas</span>
            </div>
          </div>

          <div class="no-projects" v-if="!dashboard.projects || dashboard.projects.length === 0">
            <p>No hay proyectos aún. ¡Crea el primero! 🎉</p>
          </div>
        </div>

        <!-- Project Modal -->
        <div class="modal-overlay" v-if="showProjectModal" @click.self="showProjectModal = false">
          <div class="modal">
            <div class="modal-header">
              <h3>{{ editingProject ? 'Editar Proyecto' : 'Nuevo Proyecto' }}</h3>
              <button @click="showProjectModal = false" class="close-btn">✕</button>
            </div>
            <form @submit.prevent="saveProject" class="modal-form">
              <div class="form-group">
                <label>Nombre</label>
                <input v-model="projectForm.name" type="text" required placeholder="Nombre del proyecto" />
              </div>
              <div class="form-group">
                <label>Descripción</label>
                <textarea v-model="projectForm.description" placeholder="Descripción opcional"></textarea>
              </div>
              <div class="form-group">
                <label>Estado</label>
                <select v-model="projectForm.status">
                  <option value="active">Activo</option>
                  <option value="paused">Pausado</option>
                  <option value="completed">Completado</option>
                </select>
              </div>
              <div class="form-actions">
                <button type="button" @click="showProjectModal = false" class="cancel-btn">Cancelar</button>
                <button type="submit" class="save-btn">{{ editingProject ? 'Actualizar' : 'Crear' }}</button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- KANBAN TAB -->
      <div v-if="activeTab === 'kanban'" class="kanban-section">
        <div class="section-header">
          <h2>📋 Kanban de Tareas</h2>
          <button @click="openTaskModal()" class="add-btn">➕ Nueva Tarea</button>
        </div>

        <div class="kanban-board">
          <div class="kanban-column">
            <div class="column-header pending">
              <span class="column-icon">⏳</span>
              <h3>Pendientes</h3>
              <span class="task-count">{{ kanbanTasks.pending.length }}</span>
            </div>
            <div class="column-content" 
                 @dragover.prevent 
                 @drop="dropTask($event, 'pending')">
              <div v-for="task in kanbanTasks.pending" 
                   :key="task.id" 
                   class="kanban-card"
                   draggable="true"
                   @dragstart="dragTask($event, task)">
                <div class="card-header">
                  <span class="task-title">{{ task.title }}</span>
                  <div class="card-actions">
                    <button @click="openTaskModal(task)" class="small-btn">✏️</button>
                    <button @click="deleteTask(task.id)" class="small-btn danger">🗑️</button>
                  </div>
                </div>
                <p class="task-desc" v-if="task.description">{{ task.description }}</p>
                <div class="card-footer">
                  <span class="priority" :class="task.priority">{{ task.priority || 'medium' }}</span>
                  <span class="project-tag" v-if="getProjectName(task.projectId)">{{ getProjectName(task.projectId) }}</span>
                </div>
              </div>
              <div class="empty-column" v-if="kanbanTasks.pending.length === 0">
                Arrastra tareas aquí
              </div>
            </div>
          </div>

          <div class="kanban-column">
            <div class="column-header in_progress">
              <span class="column-icon">🔥</span>
              <h3>En Progreso</h3>
              <span class="task-count">{{ kanbanTasks.in_progress.length }}</span>
            </div>
            <div class="column-content"
                 @dragover.prevent 
                 @drop="dropTask($event, 'in_progress')">
              <div v-for="task in kanbanTasks.in_progress" 
                   :key="task.id" 
                   class="kanban-card"
                   draggable="true"
                   @dragstart="dragTask($event, task)">
                <div class="card-header">
                  <span class="task-title">{{ task.title }}</span>
                  <div class="card-actions">
                    <button @click="openTaskModal(task)" class="small-btn">✏️</button>
                    <button @click="deleteTask(task.id)" class="small-btn danger">🗑️</button>
                  </div>
                </div>
                <p class="task-desc" v-if="task.description">{{ task.description }}</p>
                <div class="card-footer">
                  <span class="priority" :class="task.priority">{{ task.priority || 'medium' }}</span>
                  <span class="project-tag" v-if="getProjectName(task.projectId)">{{ getProjectName(task.projectId) }}</span>
                </div>
              </div>
              <div class="empty-column" v-if="kanbanTasks.in_progress.length === 0">
                Arrastra tareas aquí
              </div>
            </div>
          </div>

          <div class="kanban-column">
            <div class="column-header paused">
              <span class="column-icon">⏸️</span>
              <h3>Pausadas</h3>
              <span class="task-count">{{ kanbanTasks.paused.length }}</span>
            </div>
            <div class="column-content"
                 @dragover.prevent 
                 @drop="dropTask($event, 'paused')">
              <div v-for="task in kanbanTasks.paused" 
                   :key="task.id" 
                   class="kanban-card"
                   draggable="true"
                   @dragstart="dragTask($event, task)">
                <div class="card-header">
                  <span class="task-title">{{ task.title }}</span>
                  <div class="card-actions">
                    <button @click="openTaskModal(task)" class="small-btn">✏️</button>
                    <button @click="deleteTask(task.id)" class="small-btn danger">🗑️</button>
                  </div>
                </div>
                <p class="task-desc" v-if="task.description">{{ task.description }}</p>
                <div class="card-footer">
                  <span class="priority" :class="task.priority">{{ task.priority || 'medium' }}</span>
                  <span class="project-tag" v-if="getProjectName(task.projectId)">{{ getProjectName(task.projectId) }}</span>
                </div>
              </div>
              <div class="empty-column" v-if="kanbanTasks.paused.length === 0">
                Arrastra tareas aquí
              </div>
            </div>
          </div>

          <div class="kanban-column">
            <div class="column-header completed">
              <span class="column-icon">✅</span>
              <h3>Completadas</h3>
              <span class="task-count">{{ kanbanTasks.completed.length }}</span>
            </div>
            <div class="column-content"
                 @dragover.prevent 
                 @drop="dropTask($event, 'completed')">
              <div v-for="task in kanbanTasks.completed" 
                   :key="task.id" 
                   class="kanban-card completed"
                   draggable="true"
                   @dragstart="dragTask($event, task)">
                <div class="card-header">
                  <span class="task-title">{{ task.title }}</span>
                  <div class="card-actions">
                    <button @click="openTaskModal(task)" class="small-btn">✏️</button>
                    <button @click="deleteTask(task.id)" class="small-btn danger">🗑️</button>
                  </div>
                </div>
                <p class="task-desc" v-if="task.description">{{ task.description }}</p>
                <div class="card-footer">
                  <span class="priority" :class="task.priority">{{ task.priority || 'medium' }}</span>
                  <span class="project-tag" v-if="getProjectName(task.projectId)">{{ getProjectName(task.projectId) }}</span>
                </div>
              </div>
              <div class="empty-column" v-if="kanbanTasks.completed.length === 0">
                Arrastra tareas aquí
              </div>
            </div>
          </div>
        </div>

        <!-- Task Modal -->
        <div class="modal-overlay" v-if="showTaskModal" @click.self="showTaskModal = false">
          <div class="modal">
            <div class="modal-header">
              <h3>{{ editingTask ? 'Editar Tarea' : 'Nueva Tarea' }}</h3>
              <button @click="showTaskModal = false" class="close-btn">✕</button>
            </div>
            <form @submit.prevent="saveTask" class="modal-form">
              <div class="form-group">
                <label>Título</label>
                <input v-model="taskForm.title" type="text" required placeholder="Título de la tarea" />
              </div>
              <div class="form-group">
                <label>Descripción</label>
                <textarea v-model="taskForm.description" placeholder="Descripción opcional"></textarea>
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>Proyecto</label>
                  <select v-model="taskForm.projectId" required>
                    <option value="">Seleccionar proyecto</option>
                    <option v-for="p in dashboard.projects" :key="p.id" :value="p.id">{{ p.name }}</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Prioridad</label>
                  <select v-model="taskForm.priority">
                    <option value="low">Baja</option>
                    <option value="medium">Media</option>
                    <option value="high">Alta</option>
                  </select>
                </div>
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>Estado</label>
                  <select v-model="taskForm.status">
                    <option value="pending">Pendiente</option>
                    <option value="in_progress">En Progreso</option>
                    <option value="paused">Pausada</option>
                    <option value="completed">Completada</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Fecha límite</label>
                  <input v-model="taskForm.dueDate" type="date" />
                </div>
              </div>
              <div class="form-actions">
                <button type="button" @click="showTaskModal = false" class="cancel-btn">Cancelar</button>
                <button type="submit" class="save-btn">{{ editingTask ? 'Actualizar' : 'Crear' }}</button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- AGENTS TAB -->
      <div v-if="activeTab === 'agents'" class="agents-section">
        <div class="section-header">
          <h2>🤖 Estado de Agentes OpenClaw</h2>
          <button @click="fetchAgentsStatus" class="refresh-btn">🔄 Actualizar</button>
        </div>
        
        <!-- Active Sessions -->
        <div class="agents-card" v-if="agentsStatus.sessions && agentsStatus.sessions.length > 0">
          <h3>Sesiones Activas</h3>
          <div class="agent-list">
            <div v-for="session in agentsStatus.sessions" :key="session.sessionId" class="agent-item active">
              <div class="agent-icon">🔵</div>
              <div class="agent-info">
                <span class="agent-name">{{ session.displayName || 'Main' }}</span>
                <span class="agent-channel">{{ session.channel }}</span>
              </div>
              <div class="agent-status">
                <span class="status-badge active">Activo</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Configured Agents -->
        <div class="agents-card">
          <h3>Agentes Configurados</h3>
          <div class="agent-list">
            <div class="agent-item">
              <div class="agent-icon">🔵</div>
              <div class="agent-info">
                <span class="agent-name">Main</span>
                <span class="agent-model">{{ agentsStatus.mainModel || 'minimax-m2.5:cloud' }}</span>
              </div>
              <div class="agent-status">
                <span class="status-badge" :class="agentsStatus.sessions?.length > 0 ? 'active' : 'inactive'">
                  {{ agentsStatus.sessions?.length > 0 ? 'Activo' : 'Inactivo' }}
                </span>
              </div>
            </div>
            <div class="agent-item">
              <div class="agent-icon">🟢</div>
              <div class="agent-info">
                <span class="agent-name">Delta (Tester)</span>
                <span class="agent-model">qwen2.5-coder:7b</span>
              </div>
              <div class="agent-status">
                <span class="status-badge inactive">Inactivo</span>
              </div>
            </div>
            <div class="agent-item">
              <div class="agent-icon">🟡</div>
              <div class="agent-info">
                <span class="agent-name">Gamma (Repair)</span>
                <span class="agent-model">minimax-m2.5:cloud</span>
              </div>
              <div class="agent-status">
                <span class="status-badge inactive">Inactivo</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Cron Jobs -->
        <div class="agents-card" v-if="agentsStatus.cronJobs">
          <h3>Cron Jobs</h3>
          <div class="cron-list">
            <div v-for="job in agentsStatus.cronJobs" :key="job.id" class="cron-item">
              <div class="cron-info">
                <span class="cron-name">{{ job.name }}</span>
                <span class="cron-schedule">{{ job.schedule }}</span>
              </div>
              <div class="cron-status">
                <span class="status-badge" :class="job.state?.lastStatus === 'ok' ? 'active' : 'error'">
                  {{ job.state?.lastStatus === 'ok' ? '✅' : '❌' }}
                </span>
                <span class="cron-last">{{ job.state?.lastRunAtMs ? new Date(job.state.lastRunAtMs).toLocaleString('es-CR', {hour: '2-digit', minute:'2-digit'}) : 'Nunca' }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Last Test Results -->
        <div class="agents-card" v-if="lastTestResults">
          <h3>Último Test Suite</h3>
          <div class="test-summary">
            <div class="test-stat passed">
              <span class="stat-value">{{ lastTestResults.passed }}</span>
              <span class="stat-label">Pasados</span>
            </div>
            <div class="test-stat failed">
              <span class="stat-value">{{ lastTestResults.failed }}</span>
              <span class="stat-label">Fallidos</span>
            </div>
            <div class="test-stat accuracy">
              <span class="stat-value">{{ lastTestResults.accuracy?.toFixed(0) }}%</span>
              <span class="stat-label">Accuracy</span>
            </div>
          </div>
          <div class="test-duration">Duración: {{ lastTestResults.duration_ms ? (lastTestResults.duration_ms / 1000).toFixed(1) + 's' : 'N/A' }}</div>
        </div>
        
        <button @click="runTestSuite" class="run-test-btn" :disabled="isRunningTest">
          {{ isRunningTest ? '⏳ Ejecutando tests...' : '🧪 Ejecutar Test Suite' }}
        </button>
      </div>

      <!-- EXCHANGE TAB -->
      <div v-if="activeTab === 'exchange'" class="exchange-section">
        <div class="exchange-header">
          <h2>💱 Tipo de Cambio USD/CRC</h2>
          <div class="current-rate" v-if="exchangeStats">
            <span class="rate-label">Precio Actual</span>
            <span class="rate-value">₡{{ exchangeStats.currentRate.toFixed(2) }}</span>
            <span class="rate-change" :class="exchangeStats.change >= 0 ? 'up' : 'down'">
              {{ exchangeStats.change >= 0 ? '↑' : '↓' }} ₡{{ Math.abs(exchangeStats.change).toFixed(2) }} ({{ exchangeStats.changePercent.toFixed(2) }}%)
            </span>
          </div>
        </div>

        <div class="exchange-stats">
          <div class="stat-box">
            <span class="stat-box-label">Mín (30 días)</span>
            <span class="stat-box-value">₡{{ exchangeStats?.minRate.toFixed(2) }}</span>
          </div>
          <div class="stat-box">
            <span class="stat-box-label">Máx (30 días)</span>
            <span class="stat-box-value">₡{{ exchangeStats?.maxRate.toFixed(2) }}</span>
          </div>
          <div class="stat-box">
            <span class="stat-box-label">Promedio</span>
            <span class="stat-box-value">₡{{ exchangeStats?.avgRate.toFixed(2) }}</span>
          </div>
        </div>

        <div class="chart-card full-width">
          <h3>📈 Historial 30 días</h3>
          <div class="chart-container large">
            <canvas ref="exchangeChart"></canvas>
          </div>
        </div>

        <div class="predictions-section" v-if="exchangePredictions.length > 0">
          <h3>🔮 Predicciones Próximos 7 días</h3>
          <div class="predictions-grid">
            <div class="prediction-day" v-for="pred in exchangePredictions" :key="pred.date">
              <span class="pred-date">{{ formatDate(pred.date) }}</span>
              <span class="pred-value">₡{{ pred.predicted.toFixed(2) }}</span>
              <span class="pred-trend" :class="pred.trend">
                {{ pred.trend === 'up' ? '📈' : pred.trend === 'down' ? '📉' : '➡️' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- NEWS TAB -->
      <div v-if="activeTab === 'news'" class="news-section">
        <h2>📰 Noticias de tus Intereses</h2>
        
        <div class="news-filters">
          <button class="filter-btn" :class="{ active: newsFilter === 'all' }" @click="newsFilter = 'all'">Todas</button>
          <button class="filter-btn" :class="{ active: newsFilter === 'cr' }" @click="newsFilter = 'cr'">🇨🇷 Costa Rica</button>
          <button class="filter-btn" :class="{ active: newsFilter === 'crypto' }" @click="newsFilter = 'crypto'">₿ Cripto</button>
          <button class="filter-btn" :class="{ active: newsFilter === 'tech' }" @click="newsFilter = 'tech'">💻 Tech</button>
        </div>

        <div class="news-grid">
          <div v-for="(item, idx) in filteredNews" :key="idx" class="news-card">
            <div class="news-category" :class="item.category">{{ getCategoryLabel(item.category) }}</div>
            <h3 class="news-title">{{ item.title }}</h3>
            <p class="news-desc">{{ item.description }}</p>
            <div class="news-footer">
              <span class="news-source">{{ item.source }}</span>
              <span class="news-date">{{ item.publishedAt }}</span>
            </div>
          </div>
          
          <div class="no-news" v-if="filteredNews.length === 0">
            <p>No hay noticias disponibles</p>
          </div>
        </div>
      </div>

      <!-- Motivational -->
      <div class="motivational-card" v-if="dashboard.motivational && activeTab === 'dashboard'">
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

// Agents state
const agentsStatus = ref({
  sessions: [],
  mainModel: 'minimax-m2.5:cloud',
  cronJobs: []
})
const lastTestResults = ref(null)
const isRunningTest = ref(false)

const dashboard = ref({})
const news = ref([])
const activeTab = ref('dashboard')
const activeProjectTab = ref(null)
const newsFilter = ref('all')
const progressChart = ref(null)
const taskChart = ref(null)
const exchangeChart = ref(null)
const exchangeStats = ref(null)
const exchangePredictions = ref([])

// Project modal
const showProjectModal = ref(false)
const editingProject = ref(null)
const projectForm = ref({ name: '', description: '', status: 'active' })

// Task modal
const showTaskModal = ref(false)
const editingTask = ref(null)
const taskForm = ref({ title: '', description: '', projectId: '', priority: 'medium', status: 'pending', dueDate: '' })

// Drag & drop
const draggedTask = ref(null)

const completedTasks = computed(() => dashboard.value.tasks?.filter(t => t.status === 'completed').length || 0)
const pendingTasks = computed(() => dashboard.value.tasks?.filter(t => t.status === 'pending').length || 0)
const inProgressTasks = computed(() => dashboard.value.tasks?.filter(t => t.status === 'in_progress').length || 0)

const activeProject = computed(() => !activeProjectTab.value ? null : dashboard.value.projects?.find(p => p.id === activeProjectTab.value))

const kanbanTasks = computed(() => ({
  pending: dashboard.value.tasks?.filter(t => t.status === 'pending') || [],
  in_progress: dashboard.value.tasks?.filter(t => t.status === 'in_progress') || [],
  paused: dashboard.value.tasks?.filter(t => t.status === 'paused') || [],
  completed: dashboard.value.tasks?.filter(t => t.status === 'completed') || []
}))

const filteredNews = computed(() => {
  if (newsFilter.value === 'all') return news.value
  return news.value.filter(n => n.category === newsFilter.value)
})

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('es-CR', { weekday: 'short', month: 'short', day: 'numeric' })
}

const getStatusIcon = (status) => ({ pending: '⏳', in_progress: '🔥', paused: '⏸️', completed: '✅' }[status] || '⏳')
const getCategoryLabel = (cat) => ({ cr: '🇨🇷 Costa Rica', crypto: '₿ Criptomonedas', tech: '💻 Tecnología' }[cat] || '📰')

const getProjectTasks = (projectId) => dashboard.value.tasks?.filter(t => t.projectId === projectId) || []
const getProjectTaskCount = (projectId) => {
  const tasks = getProjectTasks(projectId)
  const completed = tasks.filter(t => t.status === 'completed').length
  return `${completed}/${tasks.length}`
}
const getProjectCompletedCount = (projectId) => {
  return getProjectTasks(projectId).filter(t => t.status === 'completed').length
}
const getProjectName = (projectId) => {
  const project = dashboard.value.projects?.find(p => p.id === projectId)
  return project?.name || ''
}

// Project CRUD
const openProjectModal = (project = null) => {
  editingProject.value = project
  if (project) {
    projectForm.value = { name: project.name, description: project.description || '', status: project.status }
  } else {
    projectForm.value = { name: '', description: '', status: 'active' }
  }
  showProjectModal.value = true
}

const saveProject = async () => {
  try {
    if (editingProject.value) {
      await axios.put(`${API_URL}/projects/${editingProject.value.id}`, projectForm.value)
    } else {
      await axios.post(`${API_URL}/projects`, projectForm.value)
    }
    showProjectModal.value = false
    fetchDashboard()
  } catch (e) { console.error('Error:', e) }
}

const deleteProject = async (id) => {
  if (confirm('¿Estás seguro de eliminar este proyecto?')) {
    try {
      await axios.delete(`${API_URL}/projects/${id}`)
      fetchDashboard()
    } catch (e) { console.error('Error:', e) }
  }
}

// Task CRUD
const openTaskModal = (task = null) => {
  editingTask.value = task
  if (task) {
    taskForm.value = { 
      title: task.title, 
      description: task.description || '', 
      projectId: task.projectId, 
      priority: task.priority || 'medium',
      status: task.status,
      dueDate: task.dueDate || ''
    }
  } else {
    taskForm.value = { title: '', description: '', projectId: '', priority: 'medium', status: 'pending', dueDate: '' }
  }
  showTaskModal.value = true
}

const saveTask = async () => {
  try {
    if (editingTask.value) {
      await axios.put(`${API_URL}/tasks/${editingTask.value.id}`, taskForm.value)
    } else {
      await axios.post(`${API_URL}/tasks`, taskForm.value)
    }
    showTaskModal.value = false
    fetchDashboard()
  } catch (e) { console.error('Error:', e) }
}

const deleteTask = async (id) => {
  if (confirm('¿Estás seguro de eliminar esta tarea?')) {
    try {
      await axios.delete(`${API_URL}/tasks/${id}`)
      fetchDashboard()
    } catch (e) { console.error('Error:', e) }
  }
}

// Drag & Drop
const dragTask = (event, task) => {
  draggedTask.value = task
  event.dataTransfer.effectAllowed = 'move'
}

const dropTask = async (event, status) => {
  event.preventDefault()
  if (draggedTask.value && draggedTask.value.status !== status) {
    try {
      await axios.patch(`${API_URL}/tasks/${draggedTask.value.id}/status`, { status })
      fetchDashboard()
    } catch (e) { console.error('Error:', e) }
  }
  draggedTask.value = null
}

// Agents
const fetchAgentsStatus = async () => {
  try {
    // Get sessions from backend which runs openclaw
    const sessionsRes = await fetch('/api/agents/sessions')
    const text = await sessionsRes.text()
    // Check if it's valid JSON
    let data
    try {
      data = JSON.parse(text)
    } catch {
      // Not JSON, show fallback agents
      agentsStatus.value.sessions = [
        { key: "agent:main:main", displayName: "Main Agent", model: "minimax-m2.5:cloud", status: "idle", channel: "telegram" },
        { key: "agent:delta", displayName: "Delta", model: "qwen3.5:0.8b", status: "configured", channel: "none" },
        { key: "agent:gamma", displayName: "Gamma", model: "qwen3.5:0.8b", status: "configured", channel: "none" },
        { key: "agent:epsilon", displayName: "Epsilon", model: "qwen2.5-coder:7b", status: "configured", channel: "none" },
        { key: "agent:zeta", displayName: "Zeta", model: "minimax-m2.5:cloud", status: "configured", channel: "none" }
      ]
      return
    }
    
    if (Array.isArray(data) && data.length > 0) {
      agentsStatus.value.sessions = data
    } else if (data && data.sessions) {
      agentsStatus.value.sessions = data.sessions
    } else {
      // No active sessions - show configured agents
      agentsStatus.value.sessions = [
        { key: "agent:main:main", displayName: "Main Agent", model: "minimax-m2.5:cloud", status: "idle", channel: "telegram" },
        { key: "agent:delta", displayName: "Delta", model: "qwen3.5:0.8b", status: "configured", channel: "none" },
        { key: "agent:gamma", displayName: "Gamma", model: "qwen3.5:0.8b", status: "configured", channel: "none" },
        { key: "agent:epsilon", displayName: "Epsilon", model: "qwen2.5-coder:7b", status: "configured", channel: "none" },
        { key: "agent:zeta", displayName: "Zeta", model: "minimax-m2.5:cloud", status: "configured", channel: "none" }
      ]
    }
    const cronRes = await fetch('/cron/jobs.json')
    if (cronRes.ok) {
      const cronData = await cronRes.json()
      agentsStatus.value.cronJobs = (cronData.jobs || []).map(job => ({
        id: job.id,
        name: job.name,
        schedule: job.schedule?.expr || '',
        state: job.state
      }))
    }
  } catch (e) { console.error('Error fetching agents:', e) }
}

const runTestSuite = async () => {
  isRunningTest.value = true
  try {
    const res = await fetch('https://ipcheck.diamondsb.net/api/bot-detection/run-tests', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        test_cases: [
          { account: 'JK1019', type: 'BOT', expected: 'BOT' },
          { account: 'MC001', type: 'BOT', expected: 'BOT' },
          { account: 'test109', type: 'HUMAN', expected: 'HUMAN' }
        ],
        model: 'ipcheck-balanced'
      })
    })
    if (res.ok) lastTestResults.value = await res.json()
  } catch (e) { console.error('Error running tests:', e) }
  isRunningTest.value = false
}

const fetchDashboard = async () => {
  try {
    const { data } = await axios.get(`${API_URL}/dashboard`)
    dashboard.value = data
    if (!activeProjectTab.value && data.projects?.length > 0) {
      activeProjectTab.value = data.projects[0].id
    }
    await nextTick()
    renderCharts()
  } catch (e) { console.error('Error:', e) }
}

const fetchNews = async () => {
  try {
    const { data } = await axios.get(`${API_URL}/news`)
    news.value = data
  } catch (e) { console.error('Error:', e) }
}

const fetchExchangeData = async () => {
  try {
    const [stats, preds] = await Promise.all([
      axios.get(`${API_URL}/exchange/history`),
      axios.get(`${API_URL}/exchange/predictions`)
    ])
    exchangeStats.value = stats.data
    exchangePredictions.value = preds.data
    await nextTick()
    renderExchangeChart()
  } catch (e) { console.error('Error:', e) }
}

const setTaskStatus = async (taskId, status) => {
  try {
    await axios.patch(`${API_URL}/tasks/${taskId}/status`, { status })
    fetchDashboard()
  } catch (e) { console.error('Error:', e) }
}

let progressChartInstance = null
let taskChartInstance = null
let exchangeChartInstance = null

const renderCharts = () => {
  if (!dashboard.value.projects) return
  const projects = dashboard.value.projects || []
  const tasks = dashboard.value.tasks || []
  const statuses = { pending: 0, in_progress: 0, paused: 0, completed: 0 }
  tasks.forEach(t => { statuses[t.status] = (statuses[t.status] || 0) + 1 })

  if (progressChart.value) {
    if (progressChartInstance) progressChartInstance.destroy()
    progressChartInstance = new Chart(progressChart.value, {
      type: 'bar',
      data: {
        labels: projects.map(p => p.name),
        datasets: [{ label: 'Progreso %', data: projects.map(p => p.progress), backgroundColor: ['#8b5cf6', '#06b6d4', '#10b981', '#f59e0b', '#ef4444'], borderRadius: 8 }]
      },
      options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { y: { beginAtZero: true, max: 100, grid: { color: '#1f2937' } }, x: { grid: { display: false } } } }
    })
  }

  if (taskChart.value) {
    if (taskChartInstance) taskChartInstance.destroy()
    taskChartInstance = new Chart(taskChart.value, {
      type: 'doughnut',
      data: {
        labels: ['Pendientes', 'En Progreso', 'Pausadas', 'Completadas'],
        datasets: [{ data: [statuses.pending, statuses.in_progress, statuses.paused, statuses.completed], backgroundColor: ['#f59e0b', '#ef4444', '#6b7280', '#10b981'], borderWidth: 0 }]
      },
      options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { color: '#9ca3af' } } } }
    })
  }
}

const renderExchangeChart = () => {
  if (!exchangeChart.value || !exchangeStats.value?.last30Days) return
  if (exchangeChartInstance) exchangeChartInstance.destroy()
  
  const history = exchangeStats.value.last30Days
  exchangeChartInstance = new Chart(exchangeChart.value, {
    type: 'line',
    data: {
      labels: history.map(h => h.date.slice(5)),
      datasets: [
        { label: 'Tipo de Cambio', data: history.map(h => h.rate), borderColor: '#8b5cf6', backgroundColor: 'rgba(139, 92, 246, 0.1)', fill: true, tension: 0.4 },
        { label: 'Promedio', data: history.map(() => exchangeStats.value.avgRate), borderColor: '#06b6d4', borderDash: [5, 5], pointRadius: 0 }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { labels: { color: '#9ca3af' } } },
      scales: {
        y: { grid: { color: '#1f2937' }, ticks: { color: '#9ca3af', callback: v => '₡' + v } },
        x: { grid: { display: false }, ticks: { color: '#9ca3af' } }
      }
    }
  })
}

onMounted(() => {
  fetchDashboard()
  fetchNews()
})
</script>

<style>
.app { min-height: 100vh; background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 100%); color: #e4e4e7; }

.header { display: flex; justify-content: space-between; align-items: center; padding: 1rem 2rem; background: rgba(255,255,255,0.05); border-bottom: 1px solid rgba(255,255,255,0.1); }
.logo { display: flex; align-items: center; gap: 0.5rem; }
.logo-icon { font-size: 1.5rem; }
.logo h1 { font-size: 1.5rem; font-weight: 700; background: linear-gradient(90deg, #8b5cf6, #06b6d4); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.header-info { display: flex; gap: 1rem; align-items: center; }
.weather-compact { display: flex; align-items: center; gap: 0.5rem; padding: 0.5rem 1rem; background: rgba(255,255,255,0.1); border-radius: 8px; }
.weather-compact .temp { font-weight: 700; font-size: 1.1rem; }
.weather-compact .condition { color: #9ca3af; font-size: 0.9rem; }
.exchange { padding: 0.5rem 1rem; background: rgba(255,255,255,0.1); border-radius: 8px; font-size: 0.9rem; }

.main { padding: 2rem; max-width: 1400px; margin: 0 auto; }

/* Nav Tabs */
.nav-tabs { display: flex; gap: 0.5rem; margin-bottom: 2rem; flex-wrap: wrap; }
.nav-tab { padding: 0.75rem 1.5rem; background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; color: #9ca3af; cursor: pointer; font-size: 1rem; }
.nav-tab.active { background: #8b5cf6; color: white; border-color: #8b5cf6; }

/* Weather */
.weather-section { background: linear-gradient(135deg, #1e3a5f 0%, #0d1b2a 100%); border-radius: 16px; padding: 1.5rem; margin-bottom: 2rem; border: 1px solid rgba(255,255,255,0.1); }
.weather-main { display: flex; align-items: center; gap: 1.5rem; margin-bottom: 1.5rem; }
.weather-icon { font-size: 4rem; }
.weather-temp { display: flex; flex-direction: column; }
.temp-value { font-size: 3rem; font-weight: 700; line-height: 1; }
.temp-feels { color: #9ca3af; font-size: 0.9rem; }
.weather-location { display: flex; flex-direction: column; }
.location-name { font-size: 1.25rem; font-weight: 600; }
.condition-text { color: #9ca3af; }

.weather-details { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; padding: 1rem; background: rgba(0,0,0,0.2); border-radius: 12px; }
.weather-detail { display: flex; flex-direction: column; align-items: center; gap: 0.25rem; }
.detail-icon { font-size: 1.5rem; }
.detail-label { color: #9ca3af; font-size: 0.8rem; }
.detail-value { font-weight: 600; font-size: 0.95rem; }

.weather-forecast { margin-top: 1.5rem; }
.weather-forecast h4 { color: #9ca3af; margin-bottom: 1rem; font-size: 0.9rem; text-transform: uppercase; }
.forecast-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 0.75rem; }
.forecast-day { display: flex; flex-direction: column; align-items: center; gap: 0.25rem; padding: 0.75rem; background: rgba(255,255,255,0.05); border-radius: 8px; }
.forecast-date { font-size: 0.8rem; color: #9ca3af; }
.forecast-icon { font-size: 1.5rem; }
.forecast-temp { font-size: 0.85rem; font-weight: 600; }
.forecast-precip { font-size: 0.75rem; color: #06b6d4; }

/* Stats */
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
.stat-card { display: flex; align-items: center; gap: 1rem; padding: 1.5rem; background: rgba(255,255,255,0.05); border-radius: 12px; border: 1px solid rgba(255,255,255,0.1); }
.stat-icon { font-size: 2rem; }
.stat-value { display: block; font-size: 2rem; font-weight: 700; }
.stat-label { color: #9ca3af; font-size: 0.85rem; }

/* Charts */
.charts-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
.chart-card { background: rgba(255,255,255,0.05); border-radius: 12px; padding: 1.5rem; border: 1px solid rgba(255,255,255,0.1); }
.chart-card h3 { margin-bottom: 1rem; color: #e4e4e7; }
.chart-card.full-width { grid-column: 1 / -1; }
.chart-container { height: 250px; }
.chart-container.large { height: 350px; }

/* Projects Tasks */
.projects-tasks-section { background: rgba(255,255,255,0.05); border-radius: 12px; padding: 1.5rem; border: 1px solid rgba(255,255,255,0.1); margin-bottom: 2rem; }
.projects-tasks-section h2 { margin-bottom: 1rem; color: #e4e4e7; }
.project-tabs { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1.5rem; }
.project-tab { padding: 0.75rem 1rem; background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; color: #9ca3af; cursor: pointer; display: flex; align-items: center; gap: 0.5rem; }
.project-tab.active { background: #8b5cf6; color: white; }
.task-count { font-size: 0.75rem; padding: 0.2rem 0.5rem; background: rgba(0,0,0,0.2); border-radius: 10px; }
.project-tasks-content { background: rgba(0,0,0,0.2); border-radius: 8px; padding: 1.5rem; }
.tasks-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
.tasks-header h3 { color: #e4e4e7; margin: 0; }
.project-status { padding: 0.25rem 0.75rem; border-radius: 4px; font-size: 0.75rem; text-transform: uppercase; }
.project-status.active { background: #10b981; color: white; }
.project-status.paused { background: #f59e0b; color: white; }
.project-status.completed { background: #6b7280; color: white; }
.project-progress { display: flex; align-items: center; gap: 1rem; margin-bottom: 1.5rem; }
.progress-bar.large { flex: 1; height: 12px; }
.progress-bar { width: 100%; height: 6px; background: rgba(255,255,255,0.1); border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; background: linear-gradient(90deg, #8b5cf6, #06b6d4); transition: width 0.3s ease; }
.progress-text { font-size: 0.9rem; color: #9ca3af; min-width: 40px; }
.tasks-list { display: flex; flex-direction: column; gap: 0.75rem; }
.task-row { display: flex; align-items: center; gap: 1rem; padding: 1rem; background: rgba(255,255,255,0.05); border-radius: 8px; border-left: 3px solid transparent; }
.task-row.completed { border-left-color: #10b981; }
.task-row.in_progress { border-left-color: #ef4444; }
.task-row.paused { border-left-color: #f59e0b; }
.task-row.pending { border-left-color: #6b7280; }
.task-status-indicator { font-size: 1.25rem; }
.task-info { flex: 1; }
.task-title { display: block; font-weight: 500; color: #e4e4e7; }
.task-title.completed { text-decoration: line-through; color: #6b7280; }
.task-desc { display: block; font-size: 0.85rem; color: #9ca3af; }
.task-due { display: block; font-size: 0.8rem; color: #06b6d4; margin-top: 0.25rem; }
.task-actions { display: flex; gap: 0.5rem; }
.status-select { padding: 0.5rem 0.75rem; border-radius: 6px; border: 1px solid rgba(255,255,255,0.2); background: rgba(255,255,255,0.1); color: white; font-size: 0.85rem; cursor: pointer; }
.status-select.pending { border-color: #6b7280; }
.status-select.in_progress { border-color: #ef4444; }
.status-select.paused { border-color: #f59e0b; }
.status-select.completed { border-color: #10b981; }

/* Projects Section */
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; flex-wrap: wrap; gap: 1rem; }
.section-header h2 { margin: 0; }
.add-btn { padding: 0.75rem 1.5rem; background: #8b5cf6; border: none; border-radius: 8px; color: white; font-weight: 600; cursor: pointer; }
.add-btn:hover { background: #7c3aed; }

.projects-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
.project-card { background: rgba(255,255,255,0.05); border-radius: 12px; padding: 1.5rem; border: 1px solid rgba(255,255,255,0.1); }
.project-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 0.5rem; }
.project-header h3 { margin: 0; color: #e4e4e7; }
.project-actions { display: flex; gap: 0.25rem; }
.icon-btn { background: none; border: none; cursor: pointer; font-size: 1rem; padding: 0.25rem; opacity: 0.6; }
.icon-btn:hover { opacity: 1; }
.icon-btn.danger:hover { color: #ef4444; }
.project-desc { color: #9ca3af; font-size: 0.9rem; margin-bottom: 1rem; }
.project-meta { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.status-badge { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.7rem; text-transform: uppercase; }
.status-badge.active { background: #10b981; color: white; }
.status-badge.paused { background: #f59e0b; color: white; }
.status-badge.completed { background: #6b7280; color: white; }
.progress-label { font-size: 0.85rem; color: #9ca3af; }
.project-footer { display: flex; justify-content: space-between; font-size: 0.8rem; color: #6b7280; margin-top: 1rem; }
.no-projects { grid-column: 1 / -1; text-align: center; padding: 3rem; color: #6b7280; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 1rem; }
.modal { background: #1a1a2e; border-radius: 16px; width: 100%; max-width: 500px; max-height: 90vh; overflow-y: auto; border: 1px solid rgba(255,255,255,0.1); }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 1.5rem; border-bottom: 1px solid rgba(255,255,255,0.1); }
.modal-header h3 { margin: 0; }
.close-btn { background: none; border: none; color: #9ca3af; cursor: pointer; font-size: 1.25rem; }
.modal-form { padding: 1.5rem; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; color: #9ca3af; font-size: 0.9rem; }
.form-group input, .form-group select, .form-group textarea { width: 100%; padding: 0.75rem; background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; color: #e4e4e7; font-size: 1rem; }
.form-group textarea { min-height: 80px; resize: vertical; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1.5rem; }
.cancel-btn { padding: 0.75rem 1.5rem; background: rgba(255,255,255,0.1); border: none; border-radius: 8px; color: #e4e4e7; cursor: pointer; }
.save-btn { padding: 0.75rem 1.5rem; background: #8b5cf6; border: none; border-radius: 8px; color: white; font-weight: 600; cursor: pointer; }
.save-btn:hover { background: #7c3aed; }

/* Kanban */
.kanban-section { animation: fadeIn 0.3s ease; }
@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }

.kanban-board { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; overflow-x: auto; }
.kanban-column { background: rgba(255,255,255,0.03); border-radius: 12px; min-width: 250px; }
.column-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem; border-bottom: 2px solid; }
.column-header.pending { border-color: #6b7280; }
.column-header.in_progress { border-color: #ef4444; }
.column-header.paused { border-color: #f59e0b; }
.column-header.completed { border-color: #10b981; }
.column-header h3 { margin: 0; font-size: 0.95rem; flex: 1; }
.column-icon { font-size: 1.25rem; }
.task-count { background: rgba(255,255,255,0.1); padding: 0.2rem 0.5rem; border-radius: 10px; font-size: 0.8rem; }

.column-content { padding: 0.75rem; min-height: 200px; display: flex; flex-direction: column; gap: 0.5rem; }
.kanban-card { background: rgba(255,255,255,0.08); border-radius: 8px; padding: 0.75rem; cursor: grab; transition: transform 0.2s, box-shadow 0.2s; }
.kanban-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.3); }
.kanban-card:active { cursor: grabbing; }
.kanban-card.completed { opacity: 0.6; }
.kanban-card.completed .task-title { text-decoration: line-through; }
.card-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 0.5rem; }
.card-header .task-title { font-weight: 500; font-size: 0.9rem; }
.card-actions { display: flex; gap: 0.25rem; opacity: 0; transition: opacity 0.2s; }
.kanban-card:hover .card-actions { opacity: 1; }
.small-btn { background: none; border: none; cursor: pointer; font-size: 0.8rem; padding: 0.1rem; }
.small-btn.danger:hover { color: #ef4444; }
.kanban-card .task-desc { font-size: 0.8rem; color: #9ca3af; margin-bottom: 0.5rem; }
.card-footer { display: flex; justify-content: space-between; align-items: center; }
.priority { font-size: 0.7rem; padding: 0.15rem 0.4rem; border-radius: 4px; text-transform: uppercase; }
.priority.low { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.priority.medium { background: rgba(245, 158, 11, 0.2); color: #f59e0b; }
.priority.high { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
.project-tag { font-size: 0.7rem; color: #8b5cf6; }
.empty-column { display: flex; align-items: center; justify-content: center; height: 100px; border: 2px dashed rgba(255,255,255,0.1); border-radius: 8px; color: #6b7280; font-size: 0.85rem; }

/* Exchange Section */
.exchange-section { animation: fadeIn 0.3s ease; }
.exchange-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; flex-wrap: wrap; gap: 1rem; }
.exchange-header h2 { margin: 0; }
.current-rate { display: flex; align-items: baseline; gap: 1rem; }
.rate-label { color: #9ca3af; }
.rate-value { font-size: 2.5rem; font-weight: 700; color: #10b981; }
.rate-change { font-size: 1rem; padding: 0.25rem 0.75rem; border-radius: 4px; }
.rate-change.up { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
.rate-change.down { background: rgba(16, 185, 129, 0.2); color: #10b981; }

.exchange-stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin-bottom: 2rem; }
.stat-box { display: flex; flex-direction: column; align-items: center; padding: 1.5rem; background: rgba(255,255,255,0.05); border-radius: 12px; border: 1px solid rgba(255,255,255,0.1); }
.stat-box-label { color: #9ca3af; font-size: 0.85rem; margin-bottom: 0.5rem; }
.stat-box-value { font-size: 1.5rem; font-weight: 700; }

.predictions-section { margin-top: 2rem; }
.predictions-section h3 { margin-bottom: 1rem; }
.predictions-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 1rem; }
.prediction-day { display: flex; flex-direction: column; align-items: center; gap: 0.25rem; padding: 1rem; background: rgba(255,255,255,0.05); border-radius: 8px; }
.pred-date { font-size: 0.8rem; color: #9ca3af; }
.pred-value { font-size: 1.1rem; font-weight: 600; }
.pred-trend { font-size: 1.25rem; }

/* News Section */
.news-section { animation: fadeIn 0.3s ease; }
.news-section h2 { margin-bottom: 1.5rem; }
.news-filters { display: flex; gap: 0.5rem; margin-bottom: 1.5rem; flex-wrap: wrap; }
.filter-btn { padding: 0.5rem 1rem; background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; color: #9ca3af; cursor: pointer; }
.filter-btn.active { background: #8b5cf6; color: white; border-color: #8b5cf6; }

.news-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
.news-card { background: rgba(255,255,255,0.05); border-radius: 12px; padding: 1.5rem; border: 1px solid rgba(255,255,255,0.1); }
.news-category { display: inline-block; padding: 0.25rem 0.75rem; border-radius: 4px; font-size: 0.75rem; text-transform: uppercase; margin-bottom: 0.75rem; }
.news-category.cr { background: rgba(6, 182, 212, 0.2); color: #06b6d4; }
.news-category.crypto { background: rgba(245, 158, 11, 0.2); color: #f59e0b; }
.news-category.tech { background: rgba(139, 92, 246, 0.2); color: #8b5cf6; }
.news-title { font-size: 1.1rem; margin-bottom: 0.5rem; line-height: 1.4; }
.news-desc { font-size: 0.9rem; color: #9ca3af; margin-bottom: 1rem; line-height: 1.5; }
.news-footer { display: flex; justify-content: space-between; font-size: 0.8rem; color: #6b7280; }
.no-news { grid-column: 1 / -1; text-align: center; padding: 3rem; color: #6b7280; }

/* Motivational */
.motivational-card { background: linear-gradient(90deg, rgba(139, 92, 246, 0.2), rgba(6, 182, 212, 0.2)); border-radius: 12px; padding: 1.5rem; display: flex; align-items: center; gap: 1rem; margin-top: 2rem; }
.motivational-card p { font-size: 1.1rem; font-style: italic; }

/* Agents Section */
.agents-section { animation: fadeIn 0.3s ease; }
.refresh-btn { padding: 0.5rem 1rem; background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); border-radius: 8px; color: #e4e4e7; cursor: pointer; }
.refresh-btn:hover { background: rgba(255,255,255,0.2); }

.agents-card { background: rgba(255,255,255,0.05); border-radius: 12px; padding: 1.5rem; border: 1px solid rgba(255,255,255,0.1); margin-bottom: 1rem; }
.agents-card h3 { color: #e4e4e7; margin: 0 0 1rem 0; font-size: 1rem; }

.agent-list { display: flex; flex-direction: column; gap: 0.75rem; }
.agent-item { display: flex; align-items: center; gap: 1rem; padding: 1rem; background: rgba(0,0,0,0.2); border-radius: 8px; }
.agent-icon { font-size: 1.5rem; }
.agent-info { flex: 1; display: flex; flex-direction: column; }
.agent-name { font-weight: 600; color: #e4e4e7; }
.agent-model { font-size: 0.8rem; color: #9ca3af; }
.agent-channel { font-size: 0.8rem; color: #6b7280; }

.cron-list { display: flex; flex-direction: column; gap: 0.5rem; }
.cron-item { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem; background: rgba(0,0,0,0.2); border-radius: 8px; }
.cron-info { display: flex; flex-direction: column; }
.cron-name { font-weight: 500; color: #e4e4e7; }
.cron-schedule { font-size: 0.8rem; color: #9ca3af; }
.cron-status { display: flex; align-items: center; gap: 0.5rem; }
.cron-last { font-size: 0.8rem; color: #6b7280; }

.test-summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin-bottom: 1rem; }
.test-stat { display: flex; flex-direction: column; align-items: center; padding: 1rem; background: rgba(0,0,0,0.2); border-radius: 8px; }
.test-stat.passed .stat-value { color: #10b981; }
.test-stat.failed .stat-value { color: #ef4444; }
.test-stat.accuracy .stat-value { color: #8b5cf6; }
.test-stat .stat-value { font-size: 1.5rem; font-weight: 700; }
.test-stat .stat-label { font-size: 0.8rem; color: #9ca3af; }
.test-duration { font-size: 0.85rem; color: #9ca3af; text-align: center; }

.run-test-btn { width: 100%; padding: 1rem; background: linear-gradient(90deg, #8b5cf6, #06b6d4); border: none; border-radius: 8px; color: white; font-size: 1rem; font-weight: 600; cursor: pointer; margin-top: 1rem; }
.run-test-btn:hover:not(:disabled) { opacity: 0.9; }
.run-test-btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* Responsive */
@media (max-width: 768px) {
  .kanban-board { grid-template-columns: 1fr; }
  .charts-row { grid-template-columns: 1fr; }
  .form-row { grid-template-columns: 1fr; }
  .exchange-stats { grid-template-columns: 1fr; }
}
</style>
