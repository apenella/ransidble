package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockTaskExecutor struct for mocking executor
type MockTaskExecutor struct {
	mock.Mock
}

// NewMockTaskExecutor returns a new MockTaskExecutor
func NewMockTaskExecutor() *MockTaskExecutor {
	return &MockTaskExecutor{}
}

// Execute mocks the Execute method
func (m *MockTaskExecutor) Execute(task *entity.Task) error {
	args := m.Called(task)
	return args.Error(0)
}
