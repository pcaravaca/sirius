import { ref, onUnmounted } from 'vue'

const ws = ref(null)
const isConnected = ref(false)
const lastEvent = ref(null)
const listeners = new Map()

function connect(url = `ws://${window.location.host}/ws`) {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) return

  ws.value = new WebSocket(url)

  ws.value.onopen = () => {
    isConnected.value = true
    console.log('[Sirius WS] Connected')
  }

  ws.value.onclose = () => {
    isConnected.value = false
    console.log('[Sirius WS] Disconnected, reconnecting in 3s...')
    setTimeout(() => connect(url), 3000)
  }

  ws.value.onerror = () => {
    isConnected.value = false
  }

  ws.value.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      lastEvent.value = data

      // Dispatch to type-specific listeners
      const cbs = listeners.get(data.type)
      if (cbs) {
        cbs.forEach(cb => cb(data))
      }

      // Dispatch to wildcard listeners
      const wildcards = listeners.get('*')
      if (wildcards) {
        wildcards.forEach(cb => cb(data))
      }
    } catch (e) {
      console.error('[Sirius WS] Parse error:', e)
    }
  }
}

function on(type, callback) {
  if (!listeners.has(type)) {
    listeners.set(type, new Set())
  }
  listeners.get(type).add(callback)
  return () => listeners.get(type)?.delete(callback)
}

function off(type, callback) {
  listeners.get(type)?.delete(callback)
}

function send(data) {
  if (ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify(data))
  }
}

function disconnect() {
  if (ws.value) {
    ws.value.close()
    ws.value = null
  }
}

export function useWebSocket() {
  return {
    ws,
    isConnected,
    lastEvent,
    connect,
    disconnect,
    on,
    off,
    send,
  }
}