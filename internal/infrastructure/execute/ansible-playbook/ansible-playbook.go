package ansibleplaybook

import (
	"context"
	"fmt"

	model "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
)

type AnsiblePlaybookRun struct{}

func NewAnsiblePlaybookRun() *AnsiblePlaybookRun {
	return &AnsiblePlaybookRun{}
}

func (a *AnsiblePlaybookRun) Run(ctx context.Context, options *model.AnsiblePlaybookOptions) error {
	fmt.Println("[NOT IMPLEMENTED] Running Ansible playbook")
	return nil
}
