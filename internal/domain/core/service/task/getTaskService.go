package task

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

const (
	// ErrTaskNotFound represents an error when a task is not found
	ErrTaskNotFound = "task not found"
	// ErrStoreNotInitialized represents an error when the store is not initialized
	ErrStoreNotInitialized = "store not initialized"
	// ErrTaskIDNotProvided represents an error when the task id is not provided
	ErrTaskIDNotProvided = "task id not provided"
)

type GetTaskService struct {
	store  repository.TaskStorer
	logger repository.Logger
}

func NewGetTaskService(store repository.TaskStorer, logger repository.Logger) *GetTaskService {
	return &GetTaskService{
		store:  store,
		logger: logger,
	}
}

func (t *GetTaskService) GetTask(id string) (*entity.Task, error) {

	if t.store == nil {
		t.logger.Error(ErrStoreNotInitialized)
		return nil, fmt.Errorf("%s", ErrStoreNotInitialized)
	}

	if id == "" {
		t.logger.Error(ErrTaskIDNotProvided)
		return nil, fmt.Errorf("%s", ErrTaskIDNotProvided)
	}

	t.logger.Debug(fmt.Sprintf("getting task %s\n", id), map[string]interface{}{"component": "service"})

	task, err := t.store.Find(id)
	if err != nil {
		t.logger.Error("%s: %s", ErrTaskNotFound, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrTaskNotFound, err)
	}

	return task, nil
}
