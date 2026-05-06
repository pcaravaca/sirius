import { ref, onMounted } from 'vue'

export function useWebSocket() {
  const isConnected = ref(false)
  const lastMessage = ref(null)
  const listeners = {}
  let ws = null

  function connect() {
    const url = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`
    ws = new WebSocket(url)
    ws.onopen = () => isConnected.value = true
    ws.onmessage = e => { try { lastMessage.value = JSON.parse(e.data) } catch {} }
    ws.onclose = () => { isConnected.value = false; setTimeout(connect, 5000) }
  }
  function on(type, fn) { if (!listeners[type]) listeners[type] = []; listeners[type].push(fn) }
  function off(type, fn) { if (!listeners[type]) return; listeners[type] = listeners[type].filter(f => f !== fn) }

  onMounted(connect)
  return { isConnected, lastMessage, on, off, connect }
}
