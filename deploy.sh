#!/bin/bash
# deploy.sh - Sirius Command Center Deploy Script
# Usage: ./deploy.sh [build|start|stop|restart|status]
# Runs on PETERIT-PC (100.102.176.115) or locally

set -e

# ── Config ──
PROJECT_DIR="/mnt/d/Workspace/Main Repo/sirius"
BINARY="sirius-server"
PORT="${SIRIUS_PORT:-8080}"
OPENCLAW_URL="${OPENCLAW_URL:-http://localhost:18789}"
OLLAMA_URL="${OLLAMA_URL:-http://localhost:11434}"
PID_FILE="./sirius.pid"
LOG_FILE="./sirius.log"

# ── Colors ──
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${CYAN}[SIRIUS]${NC} $1"; }
ok()   { echo -e "${GREEN}[SIRIUS]${NC} ✅ $1"; }
warn() { echo -e "${YELLOW}[SIRIUS]${NC} ⚠️  $1"; }
err()  { echo -e "${RED}[SIRIUS]${NC} ❌ $1"; }

# ── Functions ──

build_frontend() {
    log "Building frontend (Vue 3 + Vite)..."
    cd "$PROJECT_DIR/web"
    
    if [ ! -d "node_modules" ]; then
        log "Installing npm dependencies..."
        npm install
    fi
    
    npm run build
    ok "Frontend built → web/dist/"
}

build_backend() {
    log "Building Go backend..."
    cd "$PROJECT_DIR"
    
    go build -o "$BINARY" ./cmd/server/
    ok "Backend built → $BINARY ($(du -h $BINARY | cut -f1))"
}

build_all() {
    log "═══ 🔨 FULL BUILD ═══"
    build_frontend
    build_backend
    ok "Build complete!"
}

is_running() {
    if [ -f "$PID_FILE" ]; then
        local pid
        pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
        rm -f "$PID_FILE"
    fi
    return 1
}

start_server() {
    if is_running; then
        warn "Already running (PID $(cat $PID_FILE))"
        return 0
    fi
    
    cd "$PROJECT_DIR"
    
    # Verify dist exists
    if [ ! -d "web/dist" ]; then
        warn "web/dist/ not found — building frontend first..."
        build_frontend
    fi
    
    # Verify binary exists
    if [ ! -f "$BINARY" ]; then
        warn "Binary not found — building backend first..."
        build_backend
    fi
    
    # Verify data dir
    mkdir -p data
    
    log "Starting Sirius on port $PORT..."
    log "  OpenClaw: $OPENCLAW_URL"
    log "  Ollama:   $OLLAMA_URL"
    
    OPENCLAW_URL="$OPENCLAW_URL" OLLAMA_URL="$OLLAMA_URL" nohup "./$BINARY" > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    
    sleep 2
    
    if is_running; then
        ok "Sirius running! PID=$(cat $PID_FILE)"
        ok "Dashboard: http://localhost:$PORT"
        ok "API:      http://localhost:$PORT/api/dashboard"
        ok "WebSocket: ws://localhost:$PORT/ws"
    else
        err "Failed to start! Check $LOG_FILE"
        tail -20 "$LOG_FILE"
        return 1
    fi
}

stop_server() {
    if ! is_running; then
        warn "Not running"
        return 0
    fi
    
    local pid
    pid=$(cat "$PID_FILE")
    log "Stopping Sirius (PID $pid)..."
    kill "$pid" 2>/dev/null || true
    
    # Wait up to 5s for graceful shutdown
    for i in $(seq 1 10); do
        if ! kill -0 "$pid" 2>/dev/null; then
            rm -f "$PID_FILE"
            ok "Stopped"
            return 0
        fi
        sleep 0.5
    done
    
    # Force kill
    warn "Force killing..."
    kill -9 "$pid" 2>/dev/null || true
    rm -f "$PID_FILE"
    ok "Force stopped"
}

show_status() {
    if is_running; then
        local pid
        pid=$(cat "$PID_FILE")
        ok "Running — PID $pid — Port $PORT"
        log "Dashboard: http://localhost:$PORT"
        log "Uptime: $(ps -o etime= -p $pid 2>/dev/null || echo 'unknown')"
        log "Memory:  $(ps -o rss= -p $pid 2>/dev/null | awk '{printf "%.1f MB\n", $1/1024}' || echo 'unknown')"
    else
        warn "Not running"
    fi
    
    # Check OpenClaw gateway
    if curl -s -o /dev/null -w "%{http_code}" "$OPENCLAW_URL/health" 2>/dev/null | grep -q "200"; then
        ok "OpenClaw gateway: ONLINE ($OPENCLAW_URL)"
    else
        warn "OpenClaw gateway: OFFLINE ($OPENCLAW_URL)"
    fi
    
    # Check Ollama
    if curl -s -o /dev/null -w "%{http_code}" "$OLLAMA_URL/api/tags" 2>/dev/null | grep -q "200"; then
        ok "Ollama: ONLINE ($OLLAMA_URL)"
    else
        warn "Ollama: OFFLINE ($OLLAMA_URL)"
    fi
}

show_logs() {
    if [ -f "$LOG_FILE" ]; then
        tail -50 "$LOG_FILE"
    else
        warn "No log file yet"
    fi
}

# ── Main ──

cd "$PROJECT_DIR" 2>/dev/null || { err "Project dir not found: $PROJECT_DIR"; exit 1; }

case "${1:-help}" in
    build)
        build_all
        ;;
    frontend|fe)
        build_frontend
        ;;
    backend|be)
        build_backend
        ;;
    start)
        start_server
        ;;
    stop)
        stop_server
        ;;
    restart)
        stop_server
        sleep 1
        start_server
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    deploy)
        build_all
        stop_server
        sleep 1
        start_server
        ;;
    *)
        echo "⭐ Sirius Command Center Deploy"
        echo ""
        echo "Usage: ./deploy.sh <command>"
        echo ""
        echo "Commands:"
        echo "  build     Build frontend + backend"
        echo "  frontend  Build Vue frontend only"
        echo "  backend   Build Go backend only"
        echo "  start     Start server"
        echo "  stop      Stop server"
        echo "  restart   Stop + start"
        echo "  deploy    Full build + restart"
        echo "  status    Show server + service status"
        echo "  logs      Show last 50 log lines"
        echo ""
        echo "Environment:"
        echo "  SIRIUS_PORT    Server port (default: 8080)"
        echo "  OPENCLAW_URL   OpenClaw gateway (default: http://localhost:18789)"
        echo "  OLLAMA_URL     Ollama server (default: http://localhost:11434)"
        ;;
esac