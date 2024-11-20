package task

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/service/executor"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository"
	persistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	// This test ensures that the GenerateID function returns a valid UUID
	t.Run("Testing the GenerateID function", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing the GenerateID function")

		s := &CreateTaskAnsiblePlaybookService{}

		id := s.GenerateID()
		assert.NotEmpty(t, id)

		_, err := uuid.Parse(id)
		assert.Nil(t, err, fmt.Sprintf("unexpected error: %v", err))
	})
}

func TestRun(t *testing.T) {
	tests := []struct {
		desc        string
		err         error
		service     *CreateTaskAnsiblePlaybookService
		task        *entity.Task
		arrangeFunc func(*testing.T, *CreateTaskAnsiblePlaybookService)
	}{
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having a nil executor",
			err:  ErrExecutorNotInitialized,
			service: NewCreateTaskAnsiblePlaybookService(
				nil,
				nil,
				nil,
				logger.NewFakeLogger(),
			),
			task:        &entity.Task{},
			arrangeFunc: nil,
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having a nil task repository",
			err:  ErrTaskRepositoryNotInitialized,
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				nil,
				nil,
				logger.NewFakeLogger(),
			),
			task:        &entity.Task{},
			arrangeFunc: nil,
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having a nil project repository",
			err:  ErrProjectRepositoryNotInitialized,
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				nil,
				logger.NewFakeLogger(),
			),
			task:        &entity.Task{},
			arrangeFunc: nil,
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having a nil task",
			err:  ErrTaskNotProvided,
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			task:        nil,
			arrangeFunc: nil,
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having a nil project id",
			err:  domainerror.NewProjectNotProvidedError(ErrProjectNotProvided),
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ProjectID: "",
			},
			arrangeFunc: nil,
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having an error on find project into the repository",
			err:  domainerror.NewProjectNotFoundError(ErrFindingProject),
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ProjectID: "project-id",
			},
			arrangeFunc: func(t *testing.T, service *CreateTaskAnsiblePlaybookService) {
				service.projectRepository.(*repository.MockProjectRepository).On("Find", "project-id").Return(nil, errors.New("error finding project"))
			},
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having an error on store the task into the repository",
			err:  fmt.Errorf("%s: %w", ErrorStoreTask, errors.New("error storing task")),
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrangeFunc: func(t *testing.T, service *CreateTaskAnsiblePlaybookService) {
				service.projectRepository.(*repository.MockProjectRepository).On("Find", "project-id").Return(&entity.Project{
					Name:      "project-id",
					Reference: "project-id",
					Format:    "plain",
					Storage:   "local",
				}, nil)
				service.taskRepository.(*persistence.MockTaskRepository).On("SafeStore", "task-id", &entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}).Return(errors.New("error storing task"))
			},
		},
		{
			desc: "Testing error running a task on the CreateTaskAnsiblePlaybookService having an error on execute the task",
			err:  fmt.Errorf("%s: %w", ErrorExecuteTask, errors.New("error executing task")),
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrangeFunc: func(t *testing.T, service *CreateTaskAnsiblePlaybookService) {
				service.projectRepository.(*repository.MockProjectRepository).On("Find", "project-id").Return(&entity.Project{
					Name:      "project-id",
					Reference: "project-id",
					Format:    "plain",
					Storage:   "local",
				}, nil)
				service.taskRepository.(*persistence.MockTaskRepository).On("SafeStore", "task-id", &entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}).Return(nil)
				service.executor.(*executor.MockExecutor).On("Execute", &entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}).Return(errors.New("error executing task"))
			},
		},
		{
			desc: "Testing success running a task on the CreateTaskAnsiblePlaybookService",
			err:  errors.New(""),
			service: NewCreateTaskAnsiblePlaybookService(
				executor.NewMockExecutor(),
				persistence.NewMockTaskRepository(),
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrangeFunc: func(t *testing.T, service *CreateTaskAnsiblePlaybookService) {
				service.projectRepository.(*repository.MockProjectRepository).On("Find", "project-id").Return(&entity.Project{
					Name:      "project-id",
					Reference: "project-id",
					Format:    "plain",
					Storage:   "local",
				}, nil)
				service.taskRepository.(*persistence.MockTaskRepository).On("SafeStore", "task-id", &entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}).Return(nil)
				service.executor.(*executor.MockExecutor).On("Execute", &entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.service)
			}

			err := test.service.Run(context.TODO(), test.task)
			if err != nil {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
