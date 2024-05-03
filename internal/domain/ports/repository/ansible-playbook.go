package repository

import (
	"context"

	model "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
)

type Runner interface {
	Run(ctx context.Context, options *model.AnsiblePlaybookOptions) error
}
