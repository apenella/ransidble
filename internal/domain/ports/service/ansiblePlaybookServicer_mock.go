package service

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockAnsiblePlaybookService struct to mock AnsiblePlaybookServicer
type MockAnsiblePlaybookService struct {
	mock.Mock
}

// NewMockAnsiblePlaybookService creates a new MockAnsiblePlaybookService
func NewMockAnsiblePlaybookService() *MockAnsiblePlaybookService {
	return &MockAnsiblePlaybookService{}
}

// GenerateID method to generate an ID
func (m *MockAnsiblePlaybookService) GenerateID() string {
	args := m.Called()
	return args.String(0)
}

// Run method to run a task
func (m *MockAnsiblePlaybookService) Run(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}
