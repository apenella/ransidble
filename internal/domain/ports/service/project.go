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

// CreateProjectServicer represents the service to create a project
type CreateProjectServicer interface {
	Create(format string, storage string, filename string, file io.Reader) error
	// Create(format string, storage string, file *multipart.FileHeader) error
}
