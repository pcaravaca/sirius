package main

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"sirius"
	"sirius/internal/config"
	"sirius/internal/handlers"
	"sirius/internal/services"
	"sirius/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	hub := websocket.NewHub()
	go hub.Run()

	projectService := services.NewProjectService()
	taskService := services.NewTaskService(cfg.TasksFile)
	weatherService := services.NewWeatherService()
	exchangeService := services.NewExchangeService()
	exchangeHistoryService := services.NewExchangeHistoryService()
	newsService := services.NewNewsService()
	agentService := services.NewAgentService(hub, cfg.OpenClawURL)

	h := handlers.NewHandler(
		projectService,
		taskService,
		weatherService,
		exchangeService,
		exchangeHistoryService,
		newsService,
		agentService,
		hub,
	)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve embedded frontend
	distFS, err := fs.Sub(sirius.FrontendFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	// Static assets from embedded FS (/assets → web/dist/assets)
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		log.Fatalf("Failed to create assets sub filesystem: %v", err)
	}
	r.StaticFS("/assets", http.FS(assetsFS))

	// SPA: serve index.html for all non-API, non-WS, non-asset routes
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || path == "/ws" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		data, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Frontend not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	api := r.Group("/api")
	{
		api.GET("/projects", h.GetProjects)
		api.POST("/projects", h.CreateProject)
		api.PUT("/projects/:id", h.UpdateProject)
		api.DELETE("/projects/:id", h.DeleteProject)

		api.GET("/tasks", h.GetTasks)
		api.POST("/tasks", h.CreateTask)
		api.PUT("/tasks/:id", h.UpdateTask)
		api.DELETE("/tasks/:id", h.DeleteTask)
		api.PATCH("/tasks/:id/toggle", h.ToggleTask)
		api.PATCH("/tasks/:id/status", h.SetTaskStatus)

		api.GET("/agents", h.GetAgents)
		api.GET("/agents/sessions", h.GetAgentSessions)
		api.POST("/agents", h.StartAgent)
		api.GET("/agents/:id", h.GetAgentOutput)
		api.POST("/agents/:id/output", h.AppendAgentOutput)
		api.POST("/agents/:id/complete", h.CompleteAgent)

		api.GET("/dashboard", h.GetDashboard)
		api.GET("/weather", h.GetWeather)
		api.GET("/exchange", h.GetExchange)
		api.GET("/exchange/history", h.GetExchangeHistory)
		api.GET("/exchange/predictions", h.GetExchangePredictions)
		api.GET("/news", h.GetNews)
	}

	r.GET("/ws", h.HandleWebSocket)

	log.Printf("🚀 Sirius v2 server running on port %s", cfg.Port)
	log.Printf("📊 Dashboard: http://localhost:%s", cfg.Port)
	r.Run(":" + cfg.Port)
}