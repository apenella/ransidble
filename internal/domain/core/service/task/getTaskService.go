package task

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

var (
	// ErrFindingTask represents an error when a task is not found
	ErrFindingTask = fmt.Errorf("error finding task")
	// ErrRepositoryNotInitialized represents an error when the store is not initialized
	ErrRepositoryNotInitialized = fmt.Errorf("task repository not initialized")
	// ErrTaskIDNotProvided represents an error when the task id is not provided
	ErrTaskIDNotProvided = fmt.Errorf("task id not provided")
)

// GetTaskService is a service to get a task
type GetTaskService struct {
	repository repository.TaskRepository
	logger     repository.Logger
}

// NewGetTaskService creates a new GetTaskService
func NewGetTaskService(repository repository.TaskRepository, logger repository.Logger) *GetTaskService {
	return &GetTaskService{
		repository: repository,
		logger:     logger,
	}
}

// GetTask returns a task by its id
func (t *GetTaskService) GetTask(id string) (*entity.Task, error) {

	if t.repository == nil {
		t.logger.Error(ErrRepositoryNotInitialized.Error(), map[string]interface{}{
			"component": "GetTaskService.GetTask",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
			"task_id":   id,
		})
		return nil, ErrRepositoryNotInitialized
	}

	if id == "" {
		t.logger.Error(ErrTaskIDNotProvided.Error(), map[string]interface{}{
			"component": "GetTaskService.GetTask",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		})

		return nil, domainerror.NewTaskNotProvidedError(ErrTaskIDNotProvided)
	}

	t.logger.Debug(fmt.Sprintf("getting task %s\n", id), map[string]interface{}{
		"component": "GetTaskService.GetTask",
		"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		"task_id":   id,
	})

	task, err := t.repository.Find(id)
	if err != nil {
		t.logger.Error("%s: %s", ErrFindingTask.Error(), err.Error(), map[string]interface{}{
			"component": "GetTaskService.GetTask",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
			"task_id":   id,
		})

		return nil, domainerror.NewTaskNotFoundError(
			fmt.Errorf("%s %s: %w", ErrFindingTask.Error(), id, err),
		)
	}

	return task, nil
}
