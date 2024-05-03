package service

import (
	"context"

	model "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
)

type AnsiblePlaybookServicer interface {
	Run(ctx context.Context, options *model.AnsiblePlaybookOptions) error
}
