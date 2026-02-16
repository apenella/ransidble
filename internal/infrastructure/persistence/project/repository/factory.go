package repository

import "github.com/apenella/ransidble/internal/domain/ports/repository"

// Factory represents the factory for project repositories
type Factory struct {
	factory map[string]repository.ProjectRepository
}

func NewFactory() *Factory {
	return &Factory{
		factory: make(map[string]repository.ProjectRepository),
	}
}

// Register registers a new project repository
func (f *Factory) Register(name string, repository repository.ProjectRepository) {
	f.factory[name] = repository
}

// Get gets a project repository by name
func (f *Factory) Get(name string) repository.ProjectRepository {
	return f.factory[name]
}
