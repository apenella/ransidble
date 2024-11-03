package executor

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/mock"
)

// MockAnsiblePlaybook represents a mock ansible playbook
type MockAnsiblePlaybook struct {
	mock.Mock
}

// NewMockAnsiblePlaybook creates a new mock ansible playbook
func NewMockAnsiblePlaybook() *MockAnsiblePlaybook {
	return &MockAnsiblePlaybook{}
}

// Run runs the mock ansible playbook
func (m *MockAnsiblePlaybook) Run(ctx context.Context, workingDir string, parameters *entity.AnsiblePlaybookParameters) error {
	args := m.Called(ctx, workingDir, parameters)
	return args.Error(0)
}
