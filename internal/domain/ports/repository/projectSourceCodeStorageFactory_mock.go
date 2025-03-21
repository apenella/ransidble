package repository

import "github.com/stretchr/testify/mock"

// MockProjectSourceCodeStorageFactory is a mock type for the SourceCodeFetchFactory
type MockProjectSourceCodeStorageFactory struct {
	mock.Mock
}

// NewMockProjectSourceCodeStorageFactory provides a mock for the SourceCodeFetchFactory
func NewMockProjectSourceCodeStorageFactory() *MockProjectSourceCodeStorageFactory {
	return &MockProjectSourceCodeStorageFactory{}
}

// Get provides a mock function with given fields: storage
func (m *MockProjectSourceCodeStorageFactory) Get(storage string) SourceCodeStorer {
	ret := m.Called(storage)

	var r0 SourceCodeStorer
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(SourceCodeStorer)
	}
	return r0
}
