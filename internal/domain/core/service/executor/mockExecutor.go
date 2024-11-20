package executor

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockExecutor struct for mocking executor
type MockExecutor struct {
	mock.Mock
}

// NewMockExecutor returns a new MockExecutor
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{}
}

// Execute mocks the Execute method
func (m *MockExecutor) Execute(task *entity.Task) error {
	args := m.Called(task)
	return args.Error(0)
}
