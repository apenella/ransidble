package mapper

import (
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/stretchr/testify/assert"
)

// TestToTaskResponse maps a task entity to a task response
func TestToTaskResponse(t *testing.T) {
	tests := []struct {
		desc     string
		task     *entity.Task
		mapper   *TaskMapper
		expected *response.TaskResponse
	}{
		{
			desc: "Testing task mapping",
			task: &entity.Task{
				Command:      "task-command",
				CompletedAt:  "task-completed-at",
				CreatedAt:    "task-created-at",
				ErrorMessage: "task-error-message",
				ExecutedAt:   "task-executed-at",
				ID:           "task-id",
				Parameters:   "task-parameters",
				ProjectID:    "task-project-id",
				Status:       "task-status",
			},
			expected: &response.TaskResponse{
				Command:      "task-command",
				CompletedAt:  "task-completed-at",
				CreatedAt:    "task-created-at",
				ErrorMessage: "task-error-message",
				ExecutedAt:   "task-executed-at",
				ID:           "task-id",
				Parameters:   "task-parameters",
				ProjectID:    "task-project-id",
				Status:       "task-status",
			},
			mapper: NewTaskMapper(),
		},
		{
			desc:     "Testing task mapping with empty task",
			task:     &entity.Task{},
			expected: &response.TaskResponse{},
			mapper:   NewTaskMapper(),
		},
		{
			desc:     "Testing task mapping with nil task",
			task:     nil,
			expected: &response.TaskResponse{},
			mapper:   NewTaskMapper(),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.ToTaskResponse(test.task)
			assert.Equal(t, test.expected, res)
		})
	}
}
