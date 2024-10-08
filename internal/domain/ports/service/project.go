package service

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// GetProjectServicer represents the service to get a task
type GetProjectServicer interface {
	GetProject(id string) (*entity.Project, error)
	GetProjectsList() ([]*entity.Project, error)
}
