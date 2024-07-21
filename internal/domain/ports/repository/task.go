package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// Executor represents an executor to run tasks
type Executor interface {
	Execute(task *entity.Task) error
}

// TaskRepository represents a repository to manage tasks
type TaskRepository interface {
	Find(id string) (*entity.Task, error)
	FindAll() ([]*entity.Task, error)
	Remove(id string) error
	SafeStore(id string, task *entity.Task) error
	Store(id string, task *entity.Task) error
	Update(id string, task *entity.Task) error
}
