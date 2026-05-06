import { ref } from 'vue'

export function useApi() {
  const data = ref([])
  const loading = ref(false)
  const error = ref(null)
  async function fetchAll(url) {
    loading.value = true; error.value = null
    try { const r = await fetch(url); if (!r.ok) throw Error(r.status); data.value = await r.json() }
    catch(e) { error.value = e.message; data.value = [] }
    finally { loading.value = false }
    return data.value
  }
  return { data, loading, error, fetchAll }
}

export function useAgents() {
  const { data, loading, error, fetchAll } = useApi()
  return { agents: data, loading, error, fetchAll: () => fetchAll('/api/agents') }
}
export function useProjects() {
  const { data, loading, error, fetchAll } = useApi()
  return { projects: data, loading, error, fetchAll: () => fetchAll('/api/projects') }
}
export function useTasks() {
  const { data, loading, error, fetchAll } = useApi()
  return { tasks: data, loading, error, fetchAll: () => fetchAll('/api/op-tasks') }
}
