package local

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepository represents a mock object of the project repository
type MockProjectRepository struct {
	mock.Mock
}

// NewMockProjectRepository creates a new mock project repository
func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{}
}

// Find mock method to find a project by ID
func (m *MockProjectRepository) Find(id string) (*entity.Project, error) {
	var project *entity.Project
	args := m.Called(id)

	if args.Get(0) == nil {
		project = nil
	} else {
		project = args.Get(0).(*entity.Project)
	}

	return project, args.Error(1)
}

// FindAll mock method to find all projects
func (m *MockProjectRepository) FindAll() ([]*entity.Project, error) {
	var projects []*entity.Project

	args := m.Called()

	if args.Get(0) == nil {
		projects = nil
	} else {
		projects = args.Get(0).([]*entity.Project)
	}

	return projects, args.Error(1)
}
