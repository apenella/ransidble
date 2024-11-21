package task

import (
	"errors"
	"fmt"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		err         error
		expected    *entity.Task
		service     *GetTaskService
		arrangeFunc func(*testing.T, *GetTaskService)
	}{
		{
			desc: "Testing getting a task on the GetTaskService",
			id:   "task-id",
			err:  errors.New(""),
			expected: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			service: NewGetTaskService(
				repository.NewMockTaskRepository(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *GetTaskService) {
				service.repository.(*repository.MockTaskRepository).On("Find", "task-id").Return(&entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}, nil)
			},
		},
		{
			desc:     "Testing error getting a task on the GetTaskService having a nil task repository",
			id:       "task-id",
			err:      ErrRepositoryNotInitialized,
			expected: nil,
			service: NewGetTaskService(
				nil,
				logger.NewFakeLogger(),
			),
			arrangeFunc: nil,
		},
		{
			desc:     "Testing error getting a task on the GetTaskService having an empty task id",
			id:       "",
			err:      domainerror.NewTaskNotProvidedError(ErrTaskIDNotProvided),
			expected: nil,
			service: NewGetTaskService(
				repository.NewMockTaskRepository(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: nil,
		},
		{
			desc: "Testing error getting a task on the GetTaskService having an error on find task into the repository",
			id:   "task-id",
			err: domainerror.NewTaskNotFoundError(
				fmt.Errorf("%s %s: %w", ErrFindingTask.Error(), "task-id", errors.New("error finding task")),
			),
			expected: nil,
			service: NewGetTaskService(
				repository.NewMockTaskRepository(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *GetTaskService) {
				service.repository.(*repository.MockTaskRepository).On("Find", "task-id").Return(nil, errors.New("error finding task"))
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

			task, err := test.service.GetTask(test.id)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, task)
			}
		})
	}
}
