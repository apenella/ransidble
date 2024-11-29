package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockProjectUnpacker is a mock type for the Unpacker
type MockProjectUnpacker struct {
	mock.Mock
}

// Unpack provides a mock function with given fields: project, workingDir
func (m *MockProjectUnpacker) Unpack(project *entity.Project, workingDir string) error {
	ret := m.Called(project, workingDir)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}
