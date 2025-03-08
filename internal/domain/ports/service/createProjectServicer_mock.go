package service

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// MockCreateProjectService struct to mock CreateProjectService
type MockCreateProjectService struct {
	mock.Mock
}

// NewMockCreateProjectService creates a new MockCreateProjectService
func NewMockCreateProjectService() *MockCreateProjectService {
	return &MockCreateProjectService{}
}

// Create method to create a project
func (m *MockCreateProjectService) Create(format string, storage string, file *multipart.FileHeader) error {
	args := m.Called(format, storage, file)
	return args.Error(0)
}
