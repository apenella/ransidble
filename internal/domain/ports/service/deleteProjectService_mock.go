package service

import "github.com/stretchr/testify/mock"

// MockDeleteProjectService struct to mock DeleteProjectService
type MockDeleteProjectService struct {
	mock.Mock
}

// NewMockDeleteProjectService creates a new MockDeleteProjectService
func NewMockDeleteProjectService() *MockDeleteProjectService {
	return &MockDeleteProjectService{}
}

// Delete method to delete a project
func (m *MockDeleteProjectService) Delete(projectID string) error {
	args := m.Called(projectID)
	return args.Error(0)
}

// DeleteVersion method to delete a project version
func (m *MockDeleteProjectService) DeleteVersion(projectID string, version string) error {
	args := m.Called(projectID, version)
	return args.Error(0)
}
