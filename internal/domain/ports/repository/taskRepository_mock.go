package repository

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository struct for mocking task repository
type MockTaskRepository struct {
	mock.Mock
}

// NewMockTaskRepository returns a new MockTaskRepository
func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{}
}

// Find mocks the Find method
func (m *MockTaskRepository) Find(id string) (*entity.Task, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Task), args.Error(1)
}

// FindAll mocks the FindAll method
func (m *MockTaskRepository) FindAll() ([]*entity.Task, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.Task), args.Error(1)
}

// Remove mocks the Remove method
func (m *MockTaskRepository) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// SafeStore mocks the SafeStore method
func (m *MockTaskRepository) SafeStore(id string, task *entity.Task) error {
	args := m.Called(id, task)
	return args.Error(0)
}

// Store mocks the Store method
func (m *MockTaskRepository) Store(id string, task *entity.Task) error {
	args := m.Called(id, task)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockTaskRepository) Update(id string, task *entity.Task) error {
	args := m.Called(id, task)
	return args.Error(0)
}
