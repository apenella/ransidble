package task

import (
	"context"
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/google/uuid"
)

var (
	// ErrExecutorNotInitialized represents an error when the executor is not initialized
	ErrExecutorNotInitialized = fmt.Errorf("executor not initialized")
	// ErrorExecuteTask represents an error when executing a task
	ErrorExecuteTask = fmt.Errorf("error executing task")
	// ErrorStoreTask represents an error when storing a task
	ErrorStoreTask = fmt.Errorf("error storing task")
	// ErrTaskRepositoryNotInitialized represents an error when the task store is not initialized
	ErrTaskRepositoryNotInitialized = fmt.Errorf("task repository not initialized")
	// ErrTaskNotProvided represents an error when the task is not provided
	ErrTaskNotProvided = fmt.Errorf("task not provided")
	// ErrProjectRepositoryNotInitialized represents an error when the project store is not initialized
	ErrProjectRepositoryNotInitialized = fmt.Errorf("project repository not initialized")
	// ErrProjectNotProvided represents an error when the project is not provided
	ErrProjectNotProvided = fmt.Errorf("project not provided")
	// ErrFindingProject represents an error when the project is not found
	ErrFindingProject = fmt.Errorf("error finding project")
	// ErrSettingUpProject represents an error when setting up a project
	ErrSettingUpProject = fmt.Errorf("error setting up project")
	// ErrGeneratingRandomString represents an error when generating a random string
	ErrGeneratingRandomString = fmt.Errorf("error generating random string")
)

// CreateTaskAnsiblePlaybookService represents the service to run an Ansible playbook
type CreateTaskAnsiblePlaybookService struct {
	executor          repository.Executor
	logger            repository.Logger
	projectRepository repository.ProjectRepository
	taskRepository    repository.TaskRepository
}

// NewCreateTaskAnsiblePlaybookService creates a new CreateTaskAnsiblePlaybookService
func NewCreateTaskAnsiblePlaybookService(executor repository.Executor, taskRepo repository.TaskRepository, projectRepo repository.ProjectRepository, logger repository.Logger) *CreateTaskAnsiblePlaybookService {
	return &CreateTaskAnsiblePlaybookService{
		executor:          executor,
		logger:            logger,
		projectRepository: projectRepo,
		taskRepository:    taskRepo,
	}
}

// GenerateID generates an ID
func (s *CreateTaskAnsiblePlaybookService) GenerateID() string {
	// TODO id generatior should be injected as a dependency
	id := uuid.New().String()
	return id
}

func (s *CreateTaskAnsiblePlaybookService) Run(ctx context.Context, projectID string, task *entity.Task) error {
	var err error
	var project *entity.Project
	// var workingDir string

	if s.executor == nil {
		s.logger.Error(ErrExecutorNotInitialized.Error(), map[string]interface{}{
			"component": "CreateTaskAnsiblePlaybookService.Run",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		})
		return ErrExecutorNotInitialized
	}

	if s.taskRepository == nil {
		s.logger.Error(ErrTaskRepositoryNotInitialized.Error(), map[string]interface{}{
			"component": "CreateTaskAnsiblePlaybookService.Run",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		})
		return ErrTaskRepositoryNotInitialized
	}

	if s.projectRepository == nil {
		s.logger.Error(ErrProjectRepositoryNotInitialized.Error(), map[string]interface{}{
			"component": "CreateTaskAnsiblePlaybookService.Run",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		})
		return ErrProjectRepositoryNotInitialized
	}

	if projectID == "" {
		s.logger.Error(ErrProjectNotProvided.Error(), map[string]interface{}{
			"component": "CreateTaskAnsiblePlaybookService.Run",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		})

		return domainerror.NewProjectNotProvidedError(ErrProjectNotProvided)
	}

	if task == nil {
		s.logger.Error(ErrTaskNotProvided.Error(), map[string]interface{}{
			"component": "CreateTaskAnsiblePlaybookService.Run",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
		})
		return ErrTaskNotProvided
	}

	err = s.taskRepository.SafeStore(task.ID, task)
	if err != nil {
		s.logger.Error("%s: %s", ErrorStoreTask, err.Error(), map[string]interface{}{
			"component":  "CreateTaskAnsiblePlaybookService.Run",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/task",
			"project_id": projectID,
			"task_id":    task.ID,
		})
		return fmt.Errorf("%s: %w", ErrorStoreTask, err)
	}

	project, err = s.projectRepository.Find(projectID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrFindingProject, err.Error()), map[string]interface{}{
			"component":  "CreateTaskAnsiblePlaybookService.Run",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/task",
			"project_id": projectID,
			"task_id":    task.ID,
		})

		return domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s %s: %w", ErrFindingProject, projectID, err),
		)
	}

	task.Project = project

	// // TODO inject random generator as a dependency
	// // Generate 5 random bytes for a 10 characters length string. Each byte is 2 characters length
	// randBytes := make([]byte, 5)
	// _, err = rand.Read(randBytes)
	// if err != nil {
	// 	errorMsg := fmt.Sprintf("%s: %s", ErrGeneratingRandomString, err.Error())
	// 	s.logger.Error(errorMsg.Error())
	// 	return fmt.Errorf("%s", errorMsg)
	// }
	// randStr := hex.EncodeToString(randBytes)

	// workingDir = fmt.Sprintf("ransidble-%s-%s", projectID, randStr)

	// err = s.projectRepository.Setup(project, workingDir)
	// if err != nil {
	// 	s.logger.Error(fmt.Sprintf("%s: %s", ErrSettingUpProject, err.Error()), map[string]interface{}{"component": "CreateTaskAnsiblePlaybookService.Run", "project_id": projectID, "task_id": task.ID})
	// 	return fmt.Errorf("%s: %w", ErrSettingUpProject, err)
	// }

	s.logger.Info(fmt.Sprintf("executing task %s", task.ID))

	err = s.executor.Execute(task)
	if err != nil {
		s.logger.Error("%s: %s", ErrorExecuteTask, err.Error(),
			map[string]interface{}{
				"component":  "CreateTaskAnsiblePlaybookService.Run",
				"package":    "github.com/apenella/ransidble/internal/domain/core/service/task",
				"project_id": projectID,
				"task_id":    task.ID,
			})
		return fmt.Errorf("%s: %w", ErrorExecuteTask, err)
	}

	return nil
}
