package ansibleplaybook

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/google/uuid"
)

// AnsiblePlaybookService represents the service to run an Ansible playbook
type AnsiblePlaybookService struct {
	executor repository.Executor
}

// NewAnsiblePlaybookService creates a new AnsiblePlaybookService
func NewAnsiblePlaybookService(executor repository.Executor) *AnsiblePlaybookService {
	return &AnsiblePlaybookService{
		executor: executor,
	}
}

// GenerateID generates an ID
func (a *AnsiblePlaybookService) GenerateID() string {
	// Generate a UUID
	id := uuid.New().String()
	return id
}

func (a *AnsiblePlaybookService) Run(ctx context.Context, task *entity.Task) error {
	err := a.executor.Execute(task)
	if err != nil {
		return err
	}

	return nil
}
