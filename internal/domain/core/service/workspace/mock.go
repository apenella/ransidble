package workspace

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/stretchr/testify/mock"
)

// MockWorkspace represents a mock workspace
type MockWorkspace struct {
	mock.Mock
}

// Prepare prepares the mock workspace
func (m *MockWorkspace) Prepare() error {
	args := m.Called()
	return args.Error(0)
}

// Cleanup cleans up the mock workspace
func (m *MockWorkspace) Cleanup() error {
	args := m.Called()
	return args.Error(0)
}

// GetWorkingDir gets the working directory of the mock workspace
func (m *MockWorkspace) GetWorkingDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// MockBuilder represents a mock builder
type MockBuilder struct {
	Workspace *MockWorkspace
}

// WithTask sets the task of the mock builder
func (m *MockBuilder) WithTask(task *entity.Task) service.WorkspaceBuilder {
	return m
}

// Build creates a new mock workspace
func (m *MockBuilder) Build() service.Workspacer {
	if m.Workspace == nil {
		m.Workspace = &MockWorkspace{}
	}

	return m.Workspace
}
