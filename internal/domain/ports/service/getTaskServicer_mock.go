package service

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockGetTaskService struct to mock GetTaskServicer
type MockGetTaskService struct {
	mock.Mock
}

// NewMockGetTaskService creates a new MockGetTaskService
func NewMockGetTaskService() *MockGetTaskService {
	return &MockGetTaskService{}
}

// GetTask method to get a task
func (m *MockGetTaskService) GetTask(id string) (*entity.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Task), args.Error(1)
}
