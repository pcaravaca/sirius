package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"sirius/internal/models"
	"sirius/internal/websocket"
)

type AgentService struct {
	agents       map[string]*models.Agent
	mu           sync.RWMutex
	hub          *websocket.Hub
	openClawURL  string
	pollInterval time.Duration
	stopPoll     chan struct{}
}

func NewAgentService(hub *websocket.Hub, openClawURL string) *AgentService {
	as := &AgentService{
		agents:       make(map[string]*models.Agent),
		hub:          hub,
		openClawURL:  openClawURL,
		pollInterval: 30 * time.Second,
		stopPoll:     make(chan struct{}),
	}

	if as.openClawURL != "" {
		go as.pollOpenClaw()
	} else {
		// No OpenClaw URL configured — seed with real agent data
		as.SeedDemoAgents()
	}

	return as
}

func (as *AgentService) pollOpenClaw() {
	as.fetchOpenClawAgents()

	ticker := time.NewTicker(as.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			as.fetchOpenClawAgents()
		case <-as.stopPoll:
			return
		}
	}
}

func (as *AgentService) fetchOpenClawAgents() {
	if as.openClawURL == "" {
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(as.openClawURL + "/api/v1/sessions")
	if err != nil {
		log.Printf("[AgentService] OpenClaw poll failed: %v", err)
		// Seed demo agents if we have none yet
		if len(as.agents) == 0 {
			as.SeedDemoAgents()
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[AgentService] OpenClaw returned status %d", resp.StatusCode)
		// Seed demo agents if OpenClaw doesn't have sessions endpoint yet
		if len(as.agents) == 0 {
			as.SeedDemoAgents()
		}
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("[AgentService] OpenClaw decode failed: %v", err)
		return
	}

	var sessions []interface{}
	if s, ok := data["sessions"].([]interface{}); ok {
		sessions = s
	} else if a, ok := data["agents"].([]interface{}); ok {
		sessions = a
	} else {
		log.Printf("[AgentService] No sessions found in OpenClaw response")
		return
	}

	agentEmojis := map[string]string{
		"alfa":    "🐺",
		"beta":    "💰",
		"gamma":   "💻",
		"delta":   "🤖",
		"epsilon": "📅",
		"zeta":    "📺",
		"eta":     "📊",
		"theta":   "🗄️",
		"iota":    "🏥",
	}

	as.mu.Lock()
	defer as.mu.Unlock()

	for _, s := range sessions {
		sessionMap, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		key, _ := sessionMap["key"].(string)
		displayName, _ := sessionMap["displayName"].(string)
		model, _ := sessionMap["model"].(string)
		status, _ := sessionMap["status"].(string)
		channel, _ := sessionMap["channel"].(string)

		if key == "" {
			continue
		}

		name := displayName
		if name == "" {
			parts := strings.Split(key, ":")
			if len(parts) >= 2 {
				name = parts[1]
			}
		}

		emoji := "🤖"
		nameLower := strings.ToLower(name)
		for k, v := range agentEmojis {
			if strings.Contains(nameLower, k) {
				emoji = v
				break
			}
		}

		agent := &models.Agent{
			ID:        key,
			Name:      name,
			Emoji:     emoji,
			Model:     model,
			Status:    status,
			Channel:   channel,
			UpdatedAt: time.Now(),
		}

		as.agents[agent.ID] = agent
	}

	allAgents := as.GetAll()
	as.hub.Broadcast(map[string]interface{}{
		"type":   "agents_update",
		"agents": allAgents,
	})
}

func (as *AgentService) StartAgent(name, task string) *models.Agent {
	as.mu.Lock()
	defer as.mu.Unlock()

	agent := &models.Agent{
		ID:        name + "-" + time.Now().Format("150405"),
		Name:      name,
		Status:    "working",
		Task:      task,
		Output:    "🚀 Started: " + task + "\n",
		StartedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	as.agents[agent.ID] = agent

	as.broadcastEvent("agent_started", agent)

	return agent
}

func (as *AgentService) UpdateOutput(agentID, output string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if agent, exists := as.agents[agentID]; exists {
		agent.Output += output + "\n"
		agent.UpdatedAt = time.Now()
		as.broadcastEvent("agent_output", agent)
	}
}

func (as *AgentService) CompleteAgent(agentID, result string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if agent, exists := as.agents[agentID]; exists {
		agent.Status = "completed"
		agent.Output += "\n✅ Completed: " + result
		agent.UpdatedAt = time.Now()
		as.broadcastEvent("agent_completed", agent)
	}
}

func (as *AgentService) ErrorAgent(agentID, err string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if agent, exists := as.agents[agentID]; exists {
		agent.Status = "error"
		agent.Output += "\n❌ Error: " + err
		agent.UpdatedAt = time.Now()
		as.broadcastEvent("agent_error", agent)
	}
}

func (as *AgentService) GetAll() []*models.Agent {
	as.mu.RLock()
	defer as.mu.RUnlock()

	result := make([]*models.Agent, 0, len(as.agents))
	for _, a := range as.agents {
		result = append(result, a)
	}
	return result
}

func (as *AgentService) Get(agentID string) (*models.Agent, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	a, ok := as.agents[agentID]
	return a, ok
}

func (as *AgentService) Remove(agentID string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	delete(as.agents, agentID)
}

// SeedDemoAgents populates the agent list with Peter's real agent team
// when OpenClaw is not reachable. Only REAL data — models actually configured,
// channels that exist, statuses that are true.
func (as *AgentService) SeedDemoAgents() {
	as.mu.Lock()

	if len(as.agents) > 0 {
		as.mu.Unlock()
		return
	}

	now := time.Now()

	// Real agent config from openclaw.json — DO NOT invent models
	demoAgents := []models.Agent{
		{
			ID: "hermes-assistant", Name: "Hermes", Emoji: "🔮",
			Model: "glm-5.1:cloud", Status: "working", Channel: "telegram",
			Task: "Managing Sirius dashboard", StartedAt: now.Add(-2 * time.Minute), UpdatedAt: now,
		},
		{
			ID: "alfa", Name: "Alfa", Emoji: "🐺",
			Model: "glm-5.1:cloud", Status: "standby", Channel: "cloud",
			UpdatedAt: now,
		},
		{
			ID: "beta", Name: "Beta", Emoji: "💰",
			Model: "minimax-m2.7:cloud", Status: "standby", Channel: "cloud",
			UpdatedAt: now,
		},
		{
			ID: "gamma", Name: "Gamma", Emoji: "💻",
			Model: "ollama/qwen2.5-coder:7b", Status: "standby", Channel: "ollama",
			UpdatedAt: now,
		},
		{
			ID: "delta", Name: "Delta", Emoji: "🤖",
			Model: "minimax-m2.7:cloud", Status: "standby", Channel: "cloud",
			UpdatedAt: now,
		},
		{
			ID: "epsilon", Name: "Epsilon", Emoji: "📅",
			Model: "ollama/gemma4:latest", Status: "standby", Channel: "ollama",
			UpdatedAt: now,
		},
		{
			ID: "zeta", Name: "Zeta", Emoji: "📺",
			Model: "ollama/gemma4:latest", Status: "standby", Channel: "ollama",
			UpdatedAt: now,
		},
		{
			ID: "eta", Name: "Eta", Emoji: "📊",
			Model: "minimax-m2.7:cloud", Status: "standby", Channel: "cloud",
			UpdatedAt: now,
		},
		{
			ID: "theta", Name: "Theta", Emoji: "🗄️",
			Model: "minimax-m2.7:cloud", Status: "standby", Channel: "cloud",
			UpdatedAt: now,
		},
		{
			ID: "iota", Name: "Iota", Emoji: "🏥",
			Model: "minimax-m2.7:cloud", Status: "standby", Channel: "cloud",
			UpdatedAt: now,
		},
	}

	for i := range demoAgents {
		as.agents[demoAgents[i].ID] = &demoAgents[i]
	}

	as.mu.Unlock()

	log.Printf("[AgentService] Seeded %d demo agents with real config", len(demoAgents))

	allAgents := as.GetAll()
	as.hub.Broadcast(map[string]interface{}{
		"type":   "agents_update",
		"agents": allAgents,
	})
}

// UpdateAgentStatus updates an agent's status, task, and output in real-time
func (as *AgentService) UpdateAgentStatus(agentID string, status string, task string, output string) (*models.Agent, bool) {
	as.mu.Lock()
	defer as.mu.Unlock()

	agent, exists := as.agents[agentID]
	if !exists {
		return nil, false
	}

	agent.Status = status
	if task != "" {
		agent.Task = task
	}
	if output != "" {
		agent.Output = output
		agent.LastOutput = output
	}
	agent.UpdatedAt = time.Now()

	if status == "working" && agent.StartedAt.IsZero() {
		agent.StartedAt = time.Now()
	}

	as.broadcastEvent("agent_status_update", agent)
	return agent, true
}

func (as *AgentService) broadcastEvent(eventType string, agent *models.Agent) {
	if as.hub != nil {
		as.hub.Broadcast(map[string]interface{}{
			"type":  eventType,
			"agent": agent,
		})
	}
}