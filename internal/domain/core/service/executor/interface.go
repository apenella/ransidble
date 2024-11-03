package executor

/*
	The AnsiblePlaybookExecutor interface is defined within the executor package because it is scoped to the executor package, rather than a core package
*/

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// AnsiblePlaybookExecutor represents the interface for the ansible playbook executor
type AnsiblePlaybookExecutor interface {
	Run(ctx context.Context, workingDir string, parameters *entity.AnsiblePlaybookParameters) error
}
