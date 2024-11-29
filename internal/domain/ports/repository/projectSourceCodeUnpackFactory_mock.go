package repository

import "github.com/stretchr/testify/mock"

// MockProjectSourceCodeUnpackFactory is a mock type for the SouceCodeUnpackFactory
type MockProjectSourceCodeUnpackFactory struct {
	mock.Mock
}

// Get provides a mock function with given fields: projectType
func (m *MockProjectSourceCodeUnpackFactory) Get(projectType string) SourceCodeUnpacker {
	ret := m.Called(projectType)

	var r0 SourceCodeUnpacker
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(SourceCodeUnpacker)
	}
	return r0
}
