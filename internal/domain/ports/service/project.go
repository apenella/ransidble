package service

import (
	"io"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// GetProjectServicer represents the service to get a task
type GetProjectServicer interface {
	GetProject(id string) (*entity.Project, error)
	GetProjectsList() ([]*entity.Project, error)
}

// CreateProjectServicer represents the service to create a project. It returns the project ID on success and an error on failure.
type CreateProjectServicer interface {
	Create(format string, storage string, filename string, file io.Reader) error
}

// DeleteProjectServicer represents the service to delete a project. It returns an error on failure.
type DeleteProjectServicer interface {
	Delete(id string) error
}
