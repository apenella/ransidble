package repository

import (
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
