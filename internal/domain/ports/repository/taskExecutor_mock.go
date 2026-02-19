package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockTaskExecutor struct for mocking executor
type MockTaskExecutor struct {
	mock.Mock
}

// Ensure MockTaskExecutor implements the Executor interface
var _ Executor = (*MockTaskExecutor)(nil)

// NewMockTaskExecutor returns a new MockTaskExecutor
func NewMockTaskExecutor() *MockTaskExecutor {
	return &MockTaskExecutor{}
}

// Execute mocks the Execute method
func (m *MockTaskExecutor) Execute(task *entity.Task) error {
	args := m.Called(task)
	return args.Error(0)
}
