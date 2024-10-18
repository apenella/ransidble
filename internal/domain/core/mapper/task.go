package mapper

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
)

// TaskMapper is responsible for mapping task entity to response
type TaskMapper struct{}

// NewTaskMapper creates a new task mapper
func NewTaskMapper() *TaskMapper {
	return &TaskMapper{}
}

// ToTaskResponse maps a task entity to a task response
func (m *TaskMapper) ToTaskResponse(task *entity.Task) *response.TaskResponse {

	return &response.TaskResponse{
		Command:      task.Command,
		CompletedAt:  task.CompletedAt,
		CreatedAt:    task.CreatedAt,
		ErrorMessage: task.ErrorMessage,
		ExecutedAt:   task.ExecutedAt,
		ID:           task.ID,
		Parameters:   task.Parameters,
		ProjectID:    task.ProjectID,
		Status:       task.Status,
	}
}
