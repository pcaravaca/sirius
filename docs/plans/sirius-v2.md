# Sirius v2 Implementation Plan

> **For Hermes:** Use subagent-driven-development skill to implement this plan task-by-task.

**Goal:** Rebuild Sirius as a gorgeous, functional real-time command center for managing AI agents and projects — dark-mode-first, mobile-first, Linear-inspired design.

**Architecture:** Vue 3 + Vite frontend consuming Go + Gin REST API + WebSocket for live updates. Single binary serves embedded SPA. OpenClaw, Ollama, and system metrics via polling + WS push. All agent types supported (OpenClaw, Claude Code, Codex, Hermes delegates).

**Tech Stack:** Go 1.24, Gin, Gorilla WebSocket, Vue 3 (Composition API), Vite 5, Chart.js, Inter font

---

## Layout (Mobile-First)

```
┌─────────────────────────────┐
│  ⭐ SIRIUS   🔔 ⚙️ 🌙28°C 💵₡452 │  ← Top bar (weather+exchange compact)
├──────┬──────────────────────┤
│ 🏠   │                      │
│ 📁   │   MAIN CONTENT       │  ← Sidebar collapses to icons on mobile
│ 🤖   │                      │
│ 📊   │                      │
│ 🖥️   │                      │
│ 💱   │                      │
│ 📰   │                      │
│ ⚙️   │                      │
├──────┴──────────────────────┤
│  🟢 Alfa working... 3m ago  │  ← Bottom status bar (live agent feed)
└─────────────────────────────┘
```

## Pages

### 1. Dashboard (🏠)
- 4 stat cards: Projects | Active Agents | Tasks Done | System Health
- Mini chart: tokens used last 24h
- Timeline: last 10 events across all sources (agent started, deploy done, task completed)
- Quick actions: Start Agent, New Task, Deploy

### 2. Projects (📁)
- Grid of project cards with: name, status badge, progress bar, last deploy time, git branch+commits
- Click card → Project Detail: tasks, deploys, linked agents, health endpoint
- CRUD operations via modal

### 3. Agents (🤖) — THE CORE PAGE
- Live agent list: emoji, name, model, status (idle/working/error/complete), channel, uptime
- Click agent → Agent Detail: live output stream (terminal-style), memory/ctx usage, token count
- Controls: Start Task, Stop, Restart, View Memory
- Agent types with icons: 🐺 OpenClaw | 🟣 Claude Code | 🟢 Codex | 🔵 Hermes
- WebSocket pushes: agent_started, agent_output, agent_completed, agent_error

### 4. Analytics (📊)
- Token usage per agent (bar chart)
- Cost tracking per model
- Task completion rate (line chart)
- Latency distribution

### 5. System (🖥️)
- Ollama: loaded models, VRAM usage, queue depth, GPU temp
- CPU/RAM/Disk usage gauges
- Docker containers status (if on same host)
- Uptime

### 6. Exchange (💱) — keep from v1
- Current rate, 30d chart, predictions

### 7. News (📰) — keep from v1
- News feed, category filters

### 8. Settings (⚙️)
- OpenClaw URL, Ollama URL, refresh intervals
- Agent type endpoints (Claude CLI path, Codex path)
- Theme toggle (dark/light)

---

## Design System (Linear-inspired)

### Colors
- Canvas: #08090a
- Panel: #0f1011  
- Surface: #191a1b
- Elevated: #28282c
- Primary text: #f7f8f8
- Secondary text: #d0d6e0
- Muted text: #8a8f98
- Brand accent: #5e6ad2
- Accent interactive: #7170ff
- Success: #10b981
- Error: #ef4444
- Warning: #f59e0b
- Border: rgba(255,255,255,0.08)
- Border subtle: rgba(255,255,255,0.05)

### Typography
- Font: Inter Variable (cv01, ss03)
- Display: 48px/510, tracking -1.056px
- H1: 32px/400, tracking -0.704px
- H2: 24px/400, tracking -0.288px
- H3: 20px/590, tracking -0.24px
- Body: 16px/400
- Small: 15px/400, tracking -0.165px
- Caption: 13px/400, tracking -0.13px
- Mono: JetBrains Mono 14px/400

### Components
- Buttons: ghost (rgba 0.02 bg), brand (#5e6ad2), pill (9999px radius)
- Cards: rgba(255,255,255,0.02) bg, rgba(255,255,255,0.08) border, 8px radius
- Badges: pill style, status colors
- Inputs: dark bg, subtle border, 6px radius
- Terminal: #0a0a0a bg, JetBrains Mono, green text for streaming output

---

## API Endpoints (Go Backend)

### Existing (keep)
- GET/POST/PUT/DELETE /api/projects, /api/tasks
- GET /api/dashboard, /api/weather, /api/exchange, /api/news
- WebSocket /ws

### New/Modified
- GET /api/agents — all agents with live status from AgentService
- GET /api/agents/:id — single agent with full output
- POST /api/agents/:id/start — start a task on an agent
- POST /api/agents/:id/stop — stop agent
- WS events: { type: "agent_started|agent_output|agent_completed|agent_error|agents_update|system_metrics", ... }
- GET /api/system — Ollama models, CPU/RAM, Docker status
- GET /api/analytics/tokens — token usage per agent
- GET /api/analytics/costs — cost breakdown

---

## Implementation Tasks

### Phase 0: Clean Slate
1. Archive v1 frontend (mv web web-v1)
2. Create new Vue 3 + Vite project
3. Set up design tokens (CSS custom properties from Linear system)
4. Create base layout component (sidebar + topbar + content + status bar)

### Phase 1: Layout + Navigation
5. Sidebar component (icons, collapse, active state)
6. Top bar component (logo, weather compact, exchange compact, notifications)
7. Status bar component (live agent feed ticker)
8. Router setup (8 pages)
9. Mobile responsive: sidebar → bottom nav on <768px

### Phase 2: Dashboard Page
10. Stat cards grid
11. Mini token chart
12. Event timeline component
13. Quick actions panel

### Phase 3: Agents Page (CORE)
14. Agent list component (live status, emoji, badges)
15. Agent detail panel (side drawer, output terminal)
16. Agent controls (start/stop/restart)
17. Terminal output viewer (streaming, monospace, auto-scroll)
18. WebSocket connection composable (useWebSocket)
19. Connect to live data from AgentService

### Phase 4: Projects Page
20. Project cards grid
21. Project detail view (tasks, agents, deploys)
22. CRUD modals

### Phase 5: System Page
23. Ollama monitor (models, VRAM, GPU)
24. System metrics (CPU, RAM, Disk gauges)
25. Docker status

### Phase 6: Analytics Page
26. Token usage chart (per agent, bar)
27. Cost tracking
28. Task completion chart (line)

### Phase 7: Backend Updates
29. Add /api/system endpoint (Ollama + system stats)
30. Add /api/analytics/tokens endpoint
31. Multi-agent support in AgentService (not just OpenClaw)
32. Git integration (last commit, branch) for projects
33. Embed frontend in Go binary (embed.FS)

### Phase 8: Polish + Deploy
34. Mobile testing (Tab A11 Plus viewport)
35. Loading states, skeleton screens
36. Error boundaries
37. Deploy script update for v2
38. Deploy to PETERIT-PC + portproxy