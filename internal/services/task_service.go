package services

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"sirius/internal/models"
)

type TaskService struct {
	tasks    map[string]*models.Task
	mu       sync.RWMutex
	dataPath string
}

func NewTaskService(dataPath string) *TaskService {
	ts := &TaskService{
		tasks:    make(map[string]*models.Task),
		dataPath: dataPath,
	}
	ts.loadFromFile()
	return ts
}

func (ts *TaskService) loadFromFile() {
	data, err := os.ReadFile(ts.dataPath)
	if err != nil {
		ts.initializeDefaults()
		ts.saveToFile()
		return
	}
	
	var tasks []*models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		ts.initializeDefaults()
		return
	}
	
	for _, t := range tasks {
		ts.tasks[t.ID] = t
	}
	
	if len(ts.tasks) == 0 {
		ts.initializeDefaults()
		ts.saveToFile()
	}
}

func (ts *TaskService) saveToFile() {
	ts.mu.RLock()
	tasks := make([]*models.Task, 0, len(ts.tasks))
	for _, t := range ts.tasks {
		tasks = append(tasks, t)
	}
	ts.mu.RUnlock()
	
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling tasks: %v\n", err)
		return
	}
	
	if err := os.WriteFile(ts.dataPath, data, 0644); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
	}
}

func (ts *TaskService) initializeDefaults() {
	defaults := []*models.Task{
		{ID: "1", ProjectID: "ipcheck-go", Title: "Terminar ipcheck-go v2.8.4", Description: "Completar desarrollo sistema geolocalización IP", Status: "pending", Priority: "high", DueDate: "2026-03-15", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", ProjectID: "ipcheck-go", Title: "Integrar detección de bots con AI", Description: "Mejorar algoritmos de detección", Status: "pending", Priority: "high", DueDate: "2026-03-10", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "3", ProjectID: "cognitrack", Title: "Terminar cognitrack", Description: "Completar monitor para servidor Ollama", Status: "pending", Priority: "high", DueDate: "2026-03-20", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "4", ProjectID: "cognitrack", Title: "Integrar métricas de GPU", Description: "Monitor de uso de GPU para Ollama", Status: "pending", Priority: "medium", DueDate: "2026-03-25", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "5", ProjectID: "hemingway-kayak", Title: "Completar módulo de escritura", Description: "Terminar features de escritura", Status: "pending", Priority: "medium", DueDate: "2026-04-01", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "6", ProjectID: "sirius", Title: "Agregar más módulos", Description: "Crear nuevos módulos para el dashboard", Status: "pending", Priority: "low", DueDate: "2026-04-15", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "7", ProjectID: "sirius", Title: "Integrar con OpenClaw", Description: "Conectar Sirius con OpenClaw", Status: "pending", Priority: "medium", DueDate: "2026-04-10", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "8", ProjectID: "Kronos", Title: "Revisar proyecto Kronos", Description: "Evaluar estado actual del proyecto", Status: "pending", Priority: "low", DueDate: "2026-04-20", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "9", ProjectID: "telegram-bot-project", Title: "Actualizar bot de Telegram", Description: "Mejoras y fixes", Status: "pending", Priority: "low", DueDate: "2026-04-25", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, t := range defaults {
		ts.tasks[t.ID] = t
	}
}

func (ts *TaskService) GetAll() []*models.Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	result := make([]*models.Task, 0, len(ts.tasks))
	for _, t := range ts.tasks {
		result = append(result, t)
	}
	return result
}

func (ts *TaskService) GetByProject(projectID string) []*models.Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	var result []*models.Task
	for _, t := range ts.tasks {
		if t.ProjectID == projectID {
			result = append(result, t)
		}
	}
	return result
}

func (ts *TaskService) Get(id string) (*models.Task, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	t, ok := ts.tasks[id]
	return t, ok
}

func (ts *TaskService) Create(task *models.Task) *models.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task.ID = generateID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	ts.tasks[task.ID] = task
	go ts.saveToFile()
	return task
}

func (ts *TaskService) Update(id string, task *models.Task) (*models.Task, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if t, exists := ts.tasks[id]; exists {
		task.ID = id
		task.UpdatedAt = time.Now()
		ts.tasks[id] = task
		go ts.saveToFile()
		return t, true
	}
	return nil, false
}

func (ts *TaskService) SetStatus(id string, status string) (*models.Task, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if t, exists := ts.tasks[id]; exists {
		t.Status = status
		t.UpdatedAt = time.Now()
		go ts.saveToFile()
		return t, true
	}
	return nil, false
}

func (ts *TaskService) Toggle(id string) (*models.Task, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if t, exists := ts.tasks[id]; exists {
		if t.Status == "completed" {
			t.Status = "pending"
		} else {
			t.Status = "completed"
		}
		t.UpdatedAt = time.Now()
		go ts.saveToFile()
		return t, true
	}
	return nil, false
}

func (ts *TaskService) Delete(id string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if _, exists := ts.tasks[id]; exists {
		delete(ts.tasks, id)
		go ts.saveToFile()
		return true
	}
	return false
}
