package service

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockGetProjectService struct to mock GetProjectService
type MockGetProjectService struct {
	mock.Mock
}

// NewMockGetProjectService creates a new MockGetProjectService
func NewMockGetProjectService() *MockGetProjectService {
	return &MockGetProjectService{}
}

// GetProject method to get a project
func (m *MockGetProjectService) GetProject(id string) (*entity.Project, error) {

	args := m.Called(id)
	return args.Get(0).(*entity.Project), args.Error(1)
}

// GetProjectsList method to get a list of projects
func (m *MockGetProjectService) GetProjectsList() ([]*entity.Project, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Project), args.Error(1)
}
