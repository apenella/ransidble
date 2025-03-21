package service

import (
	"io"

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
// func (m *MockCreateProjectService) Create(format string, storage string, file *multipart.FileHeader) error {
func (m *MockCreateProjectService) Create(format string, storage string, filename string, file io.Reader) error {
	args := m.Called(format, storage, file)
	return args.Error(0)
}
