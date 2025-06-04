package repository

import "github.com/apenella/ransidble/internal/domain/core/entity"

// backendDriverer
type backendDriverer interface {
	Read(id string) (*entity.Project, error)
	ReadAll() ([]*entity.Project, error)
	Write(id string, project *entity.Project) error
	Remove(id string) error
	Exists(id string) (bool, error)
}
