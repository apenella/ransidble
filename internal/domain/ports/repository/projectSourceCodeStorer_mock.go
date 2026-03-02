package repository

import (
	"io"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockProjectSourceCodeStorer is a mock type for the SourceCodeStorer
type MockProjectSourceCodeStorer struct {
	mock.Mock
}

// Ensure MockProjectSourceCodeStorer implements the SourceCodeStorer interface
var _ SourceCodeStorer = (*MockProjectSourceCodeStorer)(nil)

// NewMockProjectSourceCodeStorer provides a mock for the SourceCodeStorer
func NewMockProjectSourceCodeStorer() *MockProjectSourceCodeStorer {
	return &MockProjectSourceCodeStorer{}
}

// Store provides a mock function with given fields: project, file
func (m *MockProjectSourceCodeStorer) Store(project *entity.Project, file io.Reader) error {
	args := m.Called(project, file)
	return args.Error(0)
}

// Delete provides a mock function with given fields: project
func (m *MockProjectSourceCodeStorer) Delete(project *entity.Project) error {
	args := m.Called(project)
	return args.Error(0)
}
