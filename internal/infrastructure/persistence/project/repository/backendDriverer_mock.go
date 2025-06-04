package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockBackendDriverer represents a mock object of the backend driver
type MockBackendDriverer struct {
	mock.Mock
}

// NewMockBackendDriverer creates a new mock backend driver
func NewMockBackendDriverer() *MockBackendDriverer {
	return &MockBackendDriverer{}
}

// Read mock method to read a project by ID
func (m *MockBackendDriverer) Read(id string) (*entity.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Project), args.Error(1)
}

// ReadAll mock method to read all projects
func (m *MockBackendDriverer) ReadAll() ([]*entity.Project, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Project), args.Error(1)
}

// Write mock method to write a project by ID
func (m *MockBackendDriverer) Write(id string, project *entity.Project) error {
	args := m.Called(id, project)
	return args.Error(0)
}

// Remove mock method to remove a project by ID
func (m *MockBackendDriverer) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Exists mock method to check if a project exists by ID
func (m *MockBackendDriverer) Exists(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
