package repository

import "github.com/apenella/ransidble/internal/domain/core/entity"

// ProjectRepository represents a repository to manage projects
type ProjectRepository interface {
	Find(id string) (*entity.Project, error)
	FindAll() ([]*entity.Project, error)
	// Remove(id string) error
	// SafeStore(id string, project *entity.Project) error
	// Store(id string, project *entity.Project) error
	// Update(id string, project *entity.Project) error
}

// Archiver represents the component to archive and unarchive projects before executing tasks
type Archiver interface {
	Unarchive(project *entity.Project, workingDir string) error
}
