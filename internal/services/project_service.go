package services

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"sirius/internal/models"
)

type ProjectService struct {
	projects    map[string]*models.Project
	mu          sync.RWMutex
	projectsFile string
}

func NewProjectService() *ProjectService {
	ps := &ProjectService{
		projects:     make(map[string]*models.Project),
		projectsFile: getEnv("SIRIUS_PROJECTS_FILE", "data/projects.json"),
	}
	ps.loadFromFile()

	// Seed with default projects if empty
	if len(ps.projects) == 0 {
		ps.seedDefaults()
		ps.saveToFile()
	}

	return ps
}

func (ps *ProjectService) loadFromFile() {
	data, err := os.ReadFile(ps.projectsFile)
	if err != nil {
		return
	}
	var list []*models.Project
	if err := json.Unmarshal(data, &list); err != nil {
		return
	}
	for _, p := range list {
		ps.projects[p.ID] = p
	}
}

func (ps *ProjectService) saveToFile() error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	list := make([]*models.Project, 0, len(ps.projects))
	for _, p := range ps.projects {
		list = append(list, p)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ps.projectsFile, data, 0644)
}

func (ps *ProjectService) seedDefaults() {
	defaults := []*models.Project{
		{
			ID:          "ipcheck-go",
			Name:        "ipcheck-go",
			Description: "IP geolocation & bot detection AI — v2.9.411 LIVE",
			Status:      "active",
			Priority:    "high",
			Progress:    85,
			CreatedAt:   time.Now().Add(-720 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "sirius",
			Name:        "Sirius Command Center",
			Description: "Real-time dashboard for agents & projects monitoring",
			Status:      "active",
			Priority:    "high",
			Progress:    60,
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "cognitrack",
			Name:        "Cognitrack",
			Description: "Monitor para servidor Ollama con métricas de GPU",
			Status:      "active",
			Priority:    "medium",
			Progress:    30,
			CreatedAt:   time.Now().Add(-360 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "hemingway-kayak",
			Name:        "Hemingway Kayak",
			Description: "Mobile app project — APK 8.8MB compiled",
			Status:      "active",
			Priority:    "medium",
			Progress:    70,
			CreatedAt:   time.Now().Add(-240 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "accounthub",
			Name:        "AccountHub",
			Description: "Migration of Control Panel — OpenProject #37",
			Status:      "paused",
			Priority:    "high",
			Progress:    40,
			CreatedAt:   time.Now().Add(-480 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "donfeed",
			Name:        "Donfeed",
			Description: "DSB data feed service",
			Status:      "paused",
			Priority:    "low",
			Progress:    15,
			CreatedAt:   time.Now().Add(-600 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}
	for _, p := range defaults {
		ps.projects[p.ID] = p
	}
}

func (ps *ProjectService) GetAll() []*models.Project {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	result := make([]*models.Project, 0, len(ps.projects))
	for _, p := range ps.projects {
		result = append(result, p)
	}
	return result
}

func (ps *ProjectService) Get(id string) (*models.Project, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	p, ok := ps.projects[id]
	return p, ok
}

func (ps *ProjectService) Create(project *models.Project) *models.Project {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if project.ID == "" {
		project.ID = time.Now().Format("20060102150405")
	}
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	ps.projects[project.ID] = project
	ps.saveToFileLocked()

	return project
}

func (ps *ProjectService) Update(id string, project *models.Project) (*models.Project, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.projects[id]; exists {
		project.ID = id
		project.UpdatedAt = time.Now()
		ps.projects[id] = project
		ps.saveToFileLocked()
		return project, true
	}
	return nil, false
}

func (ps *ProjectService) Delete(id string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.projects[id]; exists {
		delete(ps.projects, id)
		ps.saveToFileLocked()
		return true
	}
	return false
}

// saveToFileLocked must be called with mu held
func (ps *ProjectService) saveToFileLocked() error {
	list := make([]*models.Project, 0, len(ps.projects))
	for _, p := range ps.projects {
		list = append(list, p)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ps.projectsFile, data, 0644)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}