package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/google/uuid"
)

const (
	// ErrExecutorNotInitialized represents an error when the executor is not initialized
	ErrExecutorNotInitialized = "executor not initialized"
	// ErrorExecuteTask represents an error when executing a task
	ErrorExecuteTask = "error executing task"
	// ErrorStoreTask represents an error when storing a task
	ErrorStoreTask = "error storing task"
	// ErrTaskStoreNotInitialized represents an error when the task store is not initialized
	ErrTaskStoreNotInitialized = "task store not initialized"
	// ErrTaskNotProvided represents an error when the task is not provided
	ErrTaskNotProvided = "task not provided"
)

// AnsiblePlaybookService represents the service to run an Ansible playbook
type AnsiblePlaybookService struct {
	taskStore repository.TaskStorer
	executor  repository.Executor
	logger    repository.Logger
}

// NewAnsiblePlaybookService creates a new AnsiblePlaybookService
func NewAnsiblePlaybookService(executor repository.Executor, store repository.TaskStorer, logger repository.Logger) *AnsiblePlaybookService {
	return &AnsiblePlaybookService{
		executor:  executor,
		logger:    logger,
		taskStore: store,
	}
}

// GenerateID generates an ID
func (s *AnsiblePlaybookService) GenerateID() string {
	// TODO id generatior should be injected as a dependency
	id := uuid.New().String()
	return id
}

func (s *AnsiblePlaybookService) Run(ctx context.Context, task *entity.Task) error {
	var err error

	if s.executor == nil {
		s.logger.Error(ErrExecutorNotInitialized)
		return errors.New(ErrExecutorNotInitialized)
	}

	if s.taskStore == nil {
		s.logger.Error(ErrTaskStoreNotInitialized)
		return errors.New(ErrTaskStoreNotInitialized)
	}

	if task == nil {
		s.logger.Error(ErrTaskNotProvided)
		return errors.New(ErrTaskNotProvided)
	}

	err = s.taskStore.SafeStore(task.ID, task)
	if err != nil {
		s.logger.Error("%s: %s", ErrorStoreTask, err.Error())
		return fmt.Errorf("%s: %w", ErrorStoreTask, err)
	}

	s.logger.Info("Executing task %s", task.ID)

	err = s.executor.Execute(task)
	if err != nil {
		s.logger.Error("%s: %s", ErrorExecuteTask, err.Error())
		return fmt.Errorf("%s: %w", ErrorExecuteTask, err)
	}

	return nil
}
