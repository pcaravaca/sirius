import axios from 'axios'
import { ref } from 'vue'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Request interceptor
api.interceptors.request.use((config) => {
  return config
})

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('[Sirius API]', error.message)
    return Promise.reject(error)
  }
)

// Dashboard
export function useDashboard() {
  const data = ref(null)
  const loading = ref(false)

  async function fetch() {
    loading.value = true
    try {
      const res = await api.get('/dashboard')
      data.value = res.data
    } finally {
      loading.value = false
    }
  }

  return { data, loading, fetch }
}

// Agents
export function useAgents() {
  const agents = ref([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      const res = await api.get('/agents')
      agents.value = res.data
    } finally {
      loading.value = false
    }
  }

  async function fetchOne(id) {
    const res = await api.get(`/agents/${id}`)
    return res.data
  }

  async function startAgent(name, task) {
    const res = await api.post('/agents', { name, task })
    return res.data
  }

  return { agents, loading, fetchAll, fetchOne, startAgent }
}

// Projects
export function useProjects() {
  const projects = ref([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      const res = await api.get('/projects')
      projects.value = res.data
    } finally {
      loading.value = false
    }
  }

  async function create(data) {
    const res = await api.post('/projects', data)
    return res.data
  }

  async function update(id, data) {
    const res = await api.put(`/projects/${id}`, data)
    return res.data
  }

  async function remove(id) {
    await api.delete(`/projects/${id}`)
  }

  return { projects, loading, fetchAll, create, update, remove }
}

// Tasks
export function useTasks() {
  const tasks = ref([])
  const loading = ref(false)

  async function fetchAll(projectId) {
    loading.value = true
    try {
      const params = projectId ? { projectId } : {}
      const res = await api.get('/tasks', { params })
      tasks.value = res.data
    } finally {
      loading.value = false
    }
  }

  async function create(data) {
    const res = await api.post('/tasks', data)
    return res.data
  }

  async function update(id, data) {
    const res = await api.put(`/tasks/${id}`, data)
    return res.data
  }

  async function setStatus(id, status) {
    const res = await api.patch(`/tasks/${id}/status`, { status })
    return res.data
  }

  async function remove(id) {
    await api.delete(`/tasks/${id}`)
  }

  return { tasks, loading, fetchAll, create, update, setStatus, remove }
}

// System
export function useSystem() {
  const system = ref(null)
  const loading = ref(false)

  async function fetch() {
    loading.value = true
    try {
      const res = await api.get('/system')
      system.value = res.data
    } finally {
      loading.value = false
    }
  }

  return { system, loading, fetch }
}

// Weather
export function useWeather() {
  const weather = ref(null)
  const loading = ref(false)

  async function fetch() {
    loading.value = true
    try {
      const res = await api.get('/weather')
      weather.value = res.data
    } finally {
      loading.value = false
    }
  }

  return { weather, loading, fetch }
}

// Exchange
export function useExchange() {
  const exchange = ref(null)
  const history = ref(null)
  const predictions = ref(null)
  const loading = ref(false)

  async function fetch() {
    loading.value = true
    try {
      const [rateRes, histRes, predRes] = await Promise.all([
        api.get('/exchange'),
        api.get('/exchange/history'),
        api.get('/exchange/predictions'),
      ])
      exchange.value = rateRes.data
      history.value = histRes.data
      predictions.value = predRes.data
    } finally {
      loading.value = false
    }
  }

  return { exchange, history, predictions, loading, fetch }
}

// News
export function useNews() {
  const news = ref([])
  const loading = ref(false)

  async function fetch() {
    loading.value = true
    try {
      const res = await api.get('/news')
      news.value = res.data
    } finally {
      loading.value = false
    }
  }

  return { news, loading, fetch }
}

export default api