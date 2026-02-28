package services

import (
	"sync"
	"time"
	"sirius/internal/models"
)

type TaskService struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

func NewTaskService() *TaskService {
	ts := &TaskService{
		tasks: make(map[string]*models.Task),
	}
	ts.initializeDefaults()
	return ts
}

func (ts *TaskService) initializeDefaults() {
	defaults := []*models.Task{
		{ID: "1", ProjectID: "1", Title: "Terminar ipcheck-go", Description: "Completar desarrollo v2.8.4", Status: "pending", Priority: "high", DueDate: "2026-03-15", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", ProjectID: "2", Title: "Terminar cognitrack", Description: "Completar monitor Ollama", Status: "pending", Priority: "high", DueDate: "2026-03-20", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "3", ProjectID: "3", Title: "Explorar Bitcoin", Description: "Investigar criptomonedas", Status: "pending", Priority: "medium", DueDate: "2026-04-01", CreatedAt: time.Now(), UpdatedAt: time.Now()},
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

func (ts *TaskService) Create(task *models.Task) *models.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task.ID = generateID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	ts.tasks[task.ID] = task
	return task
}

func (ts *TaskService) Update(id string, task *models.Task) (*models.Task, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if t, exists := ts.tasks[id]; exists {
		task.ID = id
		task.UpdatedAt = time.Now()
		ts.tasks[id] = task
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
		return t, true
	}
	return nil, false
}

func (ts *TaskService) Delete(id string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if _, exists := ts.tasks[id]; exists {
		delete(ts.tasks, id)
		return true
	}
	return false
}
