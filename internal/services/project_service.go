package services

import (
	"os"
	"path/filepath"
	"sync"
	"time"
	"sirius/internal/models"
)

type ProjectService struct {
	projects map[string]*models.Project
	mu       sync.RWMutex
	projectsPath string
}

func NewProjectService() *ProjectService {
	ps := &ProjectService{
		projects: make(map[string]*models.Project),
		projectsPath: os.Getenv("SIRIUS_PROJECTS_PATH"),
	}
	if ps.projectsPath == "" {
		ps.projectsPath = "/home/peter/projects"
	}
	ps.loadFromDisk()
	return ps
}

func (ps *ProjectService) loadFromDisk() {
	entries, err := os.ReadDir(ps.projectsPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip hidden dirs only
		if len(name) > 0 && name[0] == '.' {
			continue
		}

		// Check if it's a project (has go.mod, package.json, etc)
		projectPath := filepath.Join(ps.projectsPath, name)
		isProject := ps.isProjectDir(projectPath)

		status := "paused"
		progress := 0
		if isProject {
			status = "active"
			progress = ps.estimateProgress(projectPath)
		}

		description := ps.getProjectDescription(projectPath, name)

		project := &models.Project{
			ID:          name,
			Name:        name,
			Description: description,
			Status:      status,
			Priority:    "medium",
			Progress:    progress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		ps.projects[name] = project
	}
}

func (ps *ProjectService) isProjectDir(path string) bool {
	files := []string{"go.mod", "package.json", "Cargo.toml", "requirements.txt", "pyproject.toml"}
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(path, f)); err == nil {
			return true
		}
	}
	return false
}

func (ps *ProjectService) estimateProgress(path string) int {
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		return 30
	}
	return 10
}

func (ps *ProjectService) getProjectDescription(path, name string) string {
	readmeFiles := []string{"README.md", "README.txt", "readme.md"}
	for _, f := range readmeFiles {
		content, err := os.ReadFile(filepath.Join(path, f))
		if err == nil {
			for i, b := range content {
				if b == '\n' {
					return string(content[:i])
				}
			}
			return string(content)
		}
	}

	content, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err == nil {
		for i, b := range content {
			if b == '\n' {
				desc := string(content[:i])
				if len(desc) > 6 {
					return desc[6:]
				}
			}
		}
	}

	return "Proyecto en " + path
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

func (ps *ProjectService) Create(project *models.Project) *models.Project {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	project.ID = generateID()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	ps.projects[project.ID] = project
	return project
}

func (ps *ProjectService) Update(id string, project *models.Project) (*models.Project, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, exists := ps.projects[id]; exists {
		project.ID = id
		project.UpdatedAt = time.Now()
		ps.projects[id] = project
		return project, true
	}
	return nil, false
}

func (ps *ProjectService) Delete(id string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, exists := ps.projects[id]; exists {
		delete(ps.projects, id)
		return true
	}
	return false
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
