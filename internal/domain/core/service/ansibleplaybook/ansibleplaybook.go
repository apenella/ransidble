package ansibleplaybook

import (
	"context"

	model "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

type AnsiblePlaybookService struct {
	runner repository.Runner
}

func NewAnsiblePlaybookService(runner repository.Runner) *AnsiblePlaybookService {
	return &AnsiblePlaybookService{
		runner: runner,
	}
}

func (a *AnsiblePlaybookService) Run(ctx context.Context, options *model.AnsiblePlaybookOptions) error {

	err := a.runner.Run(ctx, options)
	if err != nil {
		return err
	}

	return nil
}
