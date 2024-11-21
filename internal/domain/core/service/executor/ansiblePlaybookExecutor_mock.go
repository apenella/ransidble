package executor

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockAnsiblePlaybookExecutor represents a mock ansible playbook
type MockAnsiblePlaybookExecutor struct {
	mock.Mock
}

// NewMockAnsiblePlaybookExecutor creates a new mock ansible playbook
func NewMockAnsiblePlaybookExecutor() *MockAnsiblePlaybookExecutor {
	return &MockAnsiblePlaybookExecutor{}
}

// Run runs the mock ansible playbook
func (m *MockAnsiblePlaybookExecutor) Run(ctx context.Context, workingDir string, parameters *entity.AnsiblePlaybookParameters) error {
	args := m.Called(ctx, workingDir, parameters)
	return args.Error(0)
}
