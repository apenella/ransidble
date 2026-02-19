package repository

import (
	"github.com/stretchr/testify/mock"
)

// MockWorkspace represents a mock workspace
type MockWorkspace struct {
	mock.Mock
}

// Ensure MockWorkspace implements the Workspace interface
var _ Workspacer = (*MockWorkspace)(nil)

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
