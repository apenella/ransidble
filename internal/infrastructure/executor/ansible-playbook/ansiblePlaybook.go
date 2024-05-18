package ansibleplaybook

import (
	"context"
	"strconv"

	request "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"

	"github.com/apenella/go-ansible/v2/pkg/playbook"
)

type AnsiblePlaybook struct{}

func NewAnsiblePlaybook() *AnsiblePlaybook {
	return &AnsiblePlaybook{}
}

func (a *AnsiblePlaybook) Run(ctx context.Context, parameters *request.AnsiblePlaybookParameters) error {

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{}

	if parameters.Check {
		ansiblePlaybookOptions.Check = parameters.Check
	}

	if parameters.Diff {
		ansiblePlaybookOptions.Diff = parameters.Diff
	}

	if len(parameters.ExtraVars) > 0 {
		for k, v := range parameters.ExtraVars {
			ansiblePlaybookOptions.ExtraVars[k] = v
		}
	}

	if len(parameters.ExtraVarsFile) > 0 {
		ansiblePlaybookOptions.ExtraVarsFile = append([]string{}, parameters.ExtraVarsFile...)
	}

	if parameters.FlushCache {
		ansiblePlaybookOptions.FlushCache = parameters.FlushCache
	}

	if parameters.ForceHandlers {
		ansiblePlaybookOptions.ForceHandlers = parameters.ForceHandlers
	}

	if parameters.Forks > 0 {
		ansiblePlaybookOptions.Forks = strconv.Itoa(parameters.Forks)
	}

	if len(parameters.Inventory) > 0 {
		ansiblePlaybookOptions.Inventory = parameters.Inventory
	}

	if len(parameters.Limit) > 0 {
		ansiblePlaybookOptions.Limit = parameters.Limit
	}

	if parameters.ListHosts {
		ansiblePlaybookOptions.ListHosts = parameters.ListHosts
	}

	if parameters.ListTags {
		ansiblePlaybookOptions.ListTags = parameters.ListTags
	}

	if parameters.ListTasks {
		ansiblePlaybookOptions.ListTasks = parameters.ListTasks
	}

	if len(parameters.SkipTags) > 0 {
		ansiblePlaybookOptions.SkipTags = parameters.SkipTags
	}

	if len(parameters.StartAtTask) > 0 {
		ansiblePlaybookOptions.StartAtTask = parameters.StartAtTask
	}

	if parameters.SyntaxCheck {
		ansiblePlaybookOptions.SyntaxCheck = parameters.SyntaxCheck
	}

	if len(parameters.Tags) > 0 {
		ansiblePlaybookOptions.Tags = parameters.Tags
	}

	if len(parameters.VaultID) > 0 {
		ansiblePlaybookOptions.VaultID = parameters.VaultID
	}

	if len(parameters.VaultPasswordFile) > 0 {
		ansiblePlaybookOptions.VaultPasswordFile = parameters.VaultPasswordFile
	}

	if parameters.Verbose {
		ansiblePlaybookOptions.Verbose = parameters.Verbose
	}

	if parameters.Version {
		ansiblePlaybookOptions.Version = parameters.Version
	}

	// It is temporary enabled. The idea is to remove it in the future to avoid executing playbooks into the server
	if len(parameters.Connection) > 0 {
		ansiblePlaybookOptions.Connection = parameters.Connection
	}

	if len(parameters.SCPExtraArgs) > 0 {
		ansiblePlaybookOptions.SCPExtraArgs = parameters.SCPExtraArgs
	}

	if len(parameters.SFTPExtraArgs) > 0 {
		ansiblePlaybookOptions.SFTPExtraArgs = parameters.SFTPExtraArgs
	}

	if len(parameters.SSHCommonArgs) > 0 {
		ansiblePlaybookOptions.SSHCommonArgs = parameters.SSHCommonArgs
	}

	if len(parameters.SSHExtraArgs) > 0 {
		ansiblePlaybookOptions.SSHExtraArgs = parameters.SSHExtraArgs
	}

	if parameters.Timeout > 0 {
		ansiblePlaybookOptions.Timeout = parameters.Timeout
	}

	if len(parameters.User) > 0 {
		ansiblePlaybookOptions.User = parameters.User
	}

	if parameters.Become {
		ansiblePlaybookOptions.Become = parameters.Become
	}

	if len(parameters.BecomeMethod) > 0 {
		ansiblePlaybookOptions.BecomeMethod = parameters.BecomeMethod
	}

	if len(parameters.BecomeUser) > 0 {
		ansiblePlaybookOptions.BecomeUser = parameters.BecomeUser
	}

	// Handle dependencies

	err := playbook.NewAnsiblePlaybookExecute(parameters.Playbooks...).
		WithPlaybookOptions(ansiblePlaybookOptions).
		Execute(ctx)

	if err != nil {
		return err
	}

	return nil
}
