package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/service"
)

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
