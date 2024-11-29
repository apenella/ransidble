package repository

import "github.com/stretchr/testify/mock"

// MockProjectSourceCodeFetchFactory is a mock type for the SourceCodeFetchFactory
type MockProjectSourceCodeFetchFactory struct {
	mock.Mock
}

// Get provides a mock function with given fields: projectType
func (m *MockProjectSourceCodeFetchFactory) Get(projectType string) SourceCodeFetcher {
	ret := m.Called(projectType)

	var r0 SourceCodeFetcher
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(SourceCodeFetcher)
	}
	return r0
}
