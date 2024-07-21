package repository

import "github.com/apenella/ransidble/internal/domain/core/entity"

type ProjectRepository interface {
	// Setup(project *entity.Project, workingDir string) error
	Find(id string) (*entity.Project, error)
	FindAll() ([]*entity.Project, error)
	Remove(id string) error
	SafeStore(id string, project *entity.Project) error
	Store(id string, project *entity.Project) error
	Update(id string, project *entity.Project) error
}

type Archiver interface {
	Unarchive(project *entity.Project, workingDir string) error
}
