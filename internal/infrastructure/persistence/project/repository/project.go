package repository

import (
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

// ProjectRepository represents a repository on local storage
type ProjectRepository struct {
	store  map[string]*entity.Project
	driver backendDriverer
	mutex  sync.Mutex
	logger repository.Logger
}

// NewProjectRepository returns a new ProjectRepository
func NewProjectRepository(driver backendDriverer, logger repository.Logger) *ProjectRepository {
	return &ProjectRepository{
		store:  make(map[string]*entity.Project),
		driver: driver,
		logger: logger,
	}
}

// Find finds a project by id
func (r *ProjectRepository) Find(id string) (*entity.Project, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if project, ok := r.store[id]; ok {
		return project, nil
	}

	project, err := r.driver.Read(id)
	if err != nil {
		return nil, err
	}

	r.store[id] = project

	return project, nil
}

// FindAll finds all projects
func (r *ProjectRepository) FindAll() ([]*entity.Project, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	projects := make([]*entity.Project, 0, len(r.store))
	if len(r.store) > 0 {
		for _, project := range r.store {
			projects = append(projects, project)
		}
		return projects, nil
	}

	return projects, nil
}

// Remove removes a project by id
func (r *ProjectRepository) Remove(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.store[id]; ok {
		delete(r.store, id)
	}

	return r.driver.Remove(id)
}

// SafeStore stores a project in the repository
func (r *ProjectRepository) SafeStore(id string, project *entity.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.store[id]
	if exists {
		return fmt.Errorf("project %s already exists", id)
	}

	if err := r.driver.Write(id, project); err != nil {
		return err
	}

	r.store[id] = project

	return nil
}

// Store stores a project in the repository
func (r *ProjectRepository) Store(id string, project *entity.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.driver.Write(id, project); err != nil {
		return err
	}

	r.store[id] = project

	return nil
}

// Update updates a project in the repository
func (r *ProjectRepository) Update(id string, project *entity.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.store[id]
	if !exists {
		return fmt.Errorf("project %s not found", id)
	}

	if err := r.driver.Write(id, project); err != nil {
		return err
	}

	r.store[id] = project

	return nil
}
