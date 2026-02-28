package handlers

import (
	"net/http"
	"sirius/internal/models"
	"sirius/internal/services"
	"sirius/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	projectService  *services.ProjectService
	taskService     *services.TaskService
	weatherService  *services.WeatherService
	exchangeService *services.ExchangeService
	hub             *websocket.Hub
}

func NewHandler(
	ps *services.ProjectService,
	ts *services.TaskService,
	ws *services.WeatherService,
	es *services.ExchangeService,
	hub *websocket.Hub,
) *Handler {
	return &Handler{
		projectService:  ps,
		taskService:     ts,
		weatherService:  ws,
		exchangeService: es,
		hub:             hub,
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

// Dashboard - Same data as morning report
func (h *Handler) GetDashboard(c *gin.Context) {
	weather, _ := h.weatherService.GetWeather()
	exchange, _ := h.exchangeService.GetExchange()

	dashboard := models.DashboardData{
		Projects:   convertProjects(h.projectService.GetAll()),
		Tasks:      convertTasks(h.taskService.GetAll()),
		Weather:    *weather,
		Exchange:   *exchange,
		News:       []models.NewsItem{},
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

func (h *Handler) HandleWebSocket(c *gin.Context) {
	h.hub.HandleWebSocket(c.Writer, c.Request)
}

// Helpers
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
