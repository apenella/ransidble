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
func (m *MockDeleteProjectService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
