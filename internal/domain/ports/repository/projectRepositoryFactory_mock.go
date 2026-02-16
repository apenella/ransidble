package repository

import "github.com/stretchr/testify/mock"

// MockProjectRepositoryFactory is a mock type for the ProjectRepositoryFactory
type MockProjectRepositoryFactory struct {
	mock.Mock
}

// Get provides a mock function with given fields: project repository type
func (m *MockProjectRepositoryFactory) Get(projectType string) ProjectRepository {
	ret := m.Called(projectType)

	var r0 ProjectRepository
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(ProjectRepository)
	}
	return r0
}
