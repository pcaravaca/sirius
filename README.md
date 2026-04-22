# Sirius ⭐

Dashboard modular e interactivo para gestionar proyectos y tareas en tiempo real.

## 🚀 Características

- **Dashboard en tiempo real** con gráficos interactivos (Chart.js)
- **Gestión de Proyectos** - Crear, editar, eliminar proyectos
- **Gestión de Tareas** - Kanban-style con prioridades
- **WebSocket** - Actualizaciones en tiempo real
- **Información del reporte matutino** - Clima, tipo de cambio, motivación

## 🛠️ Stack

- **Backend**: Go + Gin
- **Frontend**: Vue 3 + Vite + Chart.js
- **WebSocket**: Gorilla WebSocket

## 📦 Instalación

```bash
# Backend
cd cmd/server
go mod tidy
go run main.go

# Frontend
cd web
npm install
npm run dev
```

## 🔌 API Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/projects` | Listar proyectos |
| POST | `/api/projects` | Crear proyecto |
| PUT | `/api/projects/:id` | Actualizar proyecto |
| DELETE | `/api/projects/:id` | Eliminar proyecto |
| GET | `/api/tasks` | Listar tareas |
| POST | `/api/tasks` | Crear tarea |
| PATCH | `/api/tasks/:id/toggle` | Toggle tarea |
| GET | `/api/dashboard` | Datos completos (reporte) |
| GET | `/api/weather` | Clima |
| GET | `/api/exchange` | Tipo de cambio |

## 🌐 WebSocket

```javascript
const ws = new WebSocket('ws://localhost:8080/ws')
ws.onmessage = (e) => console.log(JSON.parse(e.data))
```

## 📁 Estructura

```
sirius/
├── cmd/server/      # Entry point Go
├── internal/
│   ├── config/      # Configuración
│   ├── handlers/    # HTTP handlers
│   ├── models/      # Modelos de datos
│   ├── services/   # Lógica de negocio
│   └── websocket/  # Hub WebSocket
└── web/             # Frontend Vue 3
```

## 🔧 Agregar Módulos

1. Crear nuevo servicio en `internal/services/`
2. Agregar endpoints en `internal/handlers/`
3. Crear componente Vue en `web/src/components/`
4. Conectar via API o WebSocket

---

🐺 Coyote Alfa - 2026
