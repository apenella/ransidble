package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockProjectSourceCodeFetcher is a mock type for the SourceCodeFetcher
type MockProjectSourceCodeFetcher struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: project, destination
func (m *MockProjectSourceCodeFetcher) Fetch(project *entity.Project, destination string) error {
	ret := m.Called(project, destination)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}
