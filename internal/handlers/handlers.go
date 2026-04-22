package handlers

import (
	"net/http"
	"sirius/internal/models"
	"sirius/internal/services"
	"sirius/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	projectService      *services.ProjectService
	taskService         *services.TaskService
	weatherService      *services.WeatherService
	exchangeService    *services.ExchangeService
	exchangeHistoryService *services.ExchangeHistoryService
	newsService        *services.NewsService
	agentService       *services.AgentService
	hub                *websocket.Hub
}

func NewHandler(
	ps *services.ProjectService,
	ts *services.TaskService,
	ws *services.WeatherService,
	es *services.ExchangeService,
	ehs *services.ExchangeHistoryService,
	ns *services.NewsService,
	as *services.AgentService,
	hub *websocket.Hub,
) *Handler {
	return &Handler{
		projectService:        ps,
		taskService:           ts,
		weatherService:        ws,
		exchangeService:       es,
		exchangeHistoryService: ehs,
		newsService:           ns,
		agentService:         as,
		hub:                   hub,
	}
}

// Project handlers
func (h *Handler) GetProjects(c *gin.Context) {
	c.JSON(http.StatusOK, h.projectService.GetAll())
}

func (h *Handler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created := h.projectService.Create(&project)
	h.hub.Broadcast(map[string]interface{}{"type": "project_created", "data": created})
	c.JSON(http.StatusCreated, created)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, ok := h.projectService.Update(id, &project)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	h.hub.Broadcast(map[string]interface{}{"type": "project_updated", "data": updated})
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if h.projectService.Delete(id) {
		h.hub.Broadcast(map[string]interface{}{"type": "project_deleted", "id": id})
		c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
}

// Task handlers
func (h *Handler) GetTasks(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID != "" {
		c.JSON(http.StatusOK, h.taskService.GetByProject(projectID))
		return
	}
	c.JSON(http.StatusOK, h.taskService.GetAll())
}

func (h *Handler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created := h.taskService.Create(&task)
	h.hub.Broadcast(map[string]interface{}{"type": "task_created", "data": created})
	c.JSON(http.StatusCreated, created)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, ok := h.taskService.Update(id, &task)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	h.hub.Broadcast(map[string]interface{}{"type": "task_updated", "data": updated})
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) SetTaskStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	validStatuses := map[string]bool{
		"pending": true, "in_progress": true, "completed": true, "paused": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	
	updated, ok := h.taskService.SetStatus(id, req.Status)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	h.hub.Broadcast(map[string]interface{}{"type": "task_status_changed", "data": updated})
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) ToggleTask(c *gin.Context) {
	id := c.Param("id")
	updated, ok := h.taskService.Toggle(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	h.hub.Broadcast(map[string]interface{}{"type": "task_toggled", "data": updated})
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if h.taskService.Delete(id) {
		h.hub.Broadcast(map[string]interface{}{"type": "task_deleted", "id": id})
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

// Agent handlers
func (h *Handler) GetAgents(c *gin.Context) {
	c.JSON(http.StatusOK, h.agentService.GetAll())
}

func (h *Handler) GetAgentSessions(c *gin.Context) {
	c.JSON(http.StatusOK, h.agentService.GetAll())
}

func (h *Handler) StartAgent(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Task string `json:"task" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	agent := h.agentService.StartAgent(req.Name, req.Task)
	h.hub.Broadcast(map[string]interface{}{
		"type": "agent_started",
		"data": agent,
	})
	c.JSON(http.StatusCreated, agent)
}

func (h *Handler) GetAgentOutput(c *gin.Context) {
	agentID := c.Param("id")
	agent, ok := h.agentService.Get(agentID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}
	c.JSON(http.StatusOK, agent)
}

func (h *Handler) AppendAgentOutput(c *gin.Context) {
	agentID := c.Param("id")
	var req struct {
		Output string `json:"output" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	h.agentService.UpdateOutput(agentID, req.Output)
	h.hub.Broadcast(map[string]interface{}{
		"type": "agent_output",
		"agentId": agentID,
		"output": req.Output,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Output appended"})
}

func (h *Handler) CompleteAgent(c *gin.Context) {
	agentID := c.Param("id")
	var req struct {
		Result string `json:"result"`
	}
	c.ShouldBindJSON(&req)
	
	h.agentService.CompleteAgent(agentID, req.Result)
	h.hub.Broadcast(map[string]interface{}{
		"type": "agent_completed",
		"agentId": agentID,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Agent completed"})
}

// Dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	weather, _ := h.weatherService.GetWeather()
	exchange, _ := h.exchangeService.GetExchange()
	news, _ := h.newsService.GetNews()

	dashboard := models.DashboardData{
		Projects:   convertProjects(h.projectService.GetAll()),
		Tasks:      convertTasks(h.taskService.GetAll()),
		Weather:    *weather,
		Exchange:   *exchange,
		News:       news,
		Motivational: "El éxito no es la meta, es el camino que recorres cada día. 🐺",
	}

	c.JSON(http.StatusOK, dashboard)
}

func (h *Handler) GetWeather(c *gin.Context) {
	weather, err := h.weatherService.GetWeather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, weather)
}

func (h *Handler) GetExchange(c *gin.Context) {
	exchange, err := h.exchangeService.GetExchange()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exchange)
}

func (h *Handler) GetExchangeHistory(c *gin.Context) {
	stats, err := h.exchangeHistoryService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetExchangePredictions(c *gin.Context) {
	predictions, err := h.exchangeHistoryService.GetPredictions(7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, predictions)
}

func (h *Handler) GetNews(c *gin.Context) {
	news, err := h.newsService.GetNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	h.hub.HandleWebSocket(c.Writer, c.Request)
}

func convertProjects(projects []*models.Project) []models.Project {
	result := make([]models.Project, len(projects))
	for i, p := range projects {
		result[i] = *p
	}
	return result
}

func convertTasks(tasks []*models.Task) []models.Task {
	result := make([]models.Task, len(tasks))
	for i, t := range tasks {
		result[i] = *t
	}
	return result
}
