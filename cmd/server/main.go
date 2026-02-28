package main

import (
	"log"
	"sirius/internal/config"
	"sirius/internal/handlers"
	"sirius/internal/services"
	"sirius/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize config
	cfg := config.Load()

	// Initialize services
	projectService := services.NewProjectService()
	taskService := services.NewTaskService()
	weatherService := services.NewWeatherService()
	exchangeService := services.NewExchangeService()

	// Initialize hub for WebSocket
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	h := handlers.NewHandler(projectService, taskService, weatherService, exchangeService, hub)

	// Setup Gin router
	r := gin.Default()
	
	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		// Projects
		api.GET("/projects", h.GetProjects)
		api.POST("/projects", h.CreateProject)
		api.PUT("/projects/:id", h.UpdateProject)
		api.DELETE("/projects/:id", h.DeleteProject)

		// Tasks
		api.GET("/tasks", h.GetTasks)
		api.POST("/tasks", h.CreateTask)
		api.PUT("/tasks/:id", h.UpdateTask)
		api.DELETE("/tasks/:id", h.DeleteTask)
		api.PATCH("/tasks/:id/toggle", h.ToggleTask)

		// Dashboard data (same as morning report)
		api.GET("/dashboard", h.GetDashboard)
		
		// Weather
		api.GET("/weather", h.GetWeather)
		
		// Exchange rate
		api.GET("/exchange", h.GetExchange)
	}

	// WebSocket
	r.GET("/ws", h.HandleWebSocket)

	// Serve static files in production
	// r.Static("/static", "./web/dist")

	log.Printf("🚀 Sirius server running on port %s", cfg.Port)
	log.Printf("📊 Dashboard: http://localhost:%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
