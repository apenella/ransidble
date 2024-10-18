package service

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// AnsiblePlaybookServicer represents the service to run an Ansible playbook
type AnsiblePlaybookServicer interface {
	GenerateID() string
	Run(ctx context.Context, task *entity.Task) error
}

// GetTaskServicer represents the service to get a task
type GetTaskServicer interface {
	GetTask(id string) (*entity.Task, error)
}
