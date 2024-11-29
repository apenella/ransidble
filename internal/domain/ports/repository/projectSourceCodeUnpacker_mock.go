package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockProjectSourceCodeUnpacker is a mock type for the SourceCodeUnpacker
type MockProjectSourceCodeUnpacker struct {
	mock.Mock
}

// Unpack provides a mock function with given fields: project, destination
func (m *MockProjectSourceCodeUnpacker) Unpack(project *entity.Project, destination string) error {
	ret := m.Called(project, destination)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}
