package service

import "github.com/apenella/ransidble/internal/domain/core/entity"

type GetTaskServicer interface {
	GetTask(id string) (*entity.Task, error)
}
