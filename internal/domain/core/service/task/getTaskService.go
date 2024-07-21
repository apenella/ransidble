package task

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

var (
	// ErrTaskNotFound represents an error when a task is not found
	ErrTaskNotFound = fmt.Errorf("task not found")
	// ErrStoreNotInitialized represents an error when the store is not initialized
	ErrStoreNotInitialized = fmt.Errorf("store not initialized")
	// ErrTaskIDNotProvided represents an error when the task id is not provided
	ErrTaskIDNotProvided = fmt.Errorf("task id not provided")
)

type GetTaskService struct {
	repository repository.TaskRepository
	logger     repository.Logger
}

func NewGetTaskService(repository repository.TaskRepository, logger repository.Logger) *GetTaskService {
	return &GetTaskService{
		repository: repository,
		logger:     logger,
	}
}

func (t *GetTaskService) GetTask(id string) (*entity.Task, error) {

	if t.repository == nil {
		t.logger.Error(ErrStoreNotInitialized.Error())
		return nil, ErrStoreNotInitialized
	}

	if id == "" {
		t.logger.Error(ErrTaskIDNotProvided.Error())
		return nil, ErrTaskIDNotProvided
	}

	t.logger.Debug(fmt.Sprintf("getting task %s\n", id), map[string]interface{}{"component": "service"})

	task, err := t.repository.Find(id)
	if err != nil {
		t.logger.Error("%s: %s", ErrTaskNotFound.Error(), err.Error())
		return nil, fmt.Errorf("%s: %w", ErrTaskNotFound.Error(), err)
	}

	return task, nil
}
