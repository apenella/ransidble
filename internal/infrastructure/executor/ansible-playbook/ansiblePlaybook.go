package ansibleplaybook

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/execute/configuration"
	"github.com/apenella/go-ansible/v2/pkg/execute/workflow"
	collection "github.com/apenella/go-ansible/v2/pkg/galaxy/collection/install"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
	"github.com/apenella/ransidble/internal/domain/core/entity"
)

const (
	// CollectionsPath represents the path where the collections are stored
	CollectionsPath = ".collections"
)

type AnsiblePlaybook struct{}

func NewAnsiblePlaybook() *AnsiblePlaybook {
	return &AnsiblePlaybook{}
}

func (a *AnsiblePlaybook) Run(ctx context.Context, workingDir string, parameters *entity.AnsiblePlaybookParameters) error {

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

	// TODO: Handle dependencies

	workflowTasks := make([]execute.Executor, 0)

	if parameters.Dependencies != nil {
		if parameters.Dependencies.Collections != nil {
			optionsFuncs := make([]collection.AnsibleGalaxyCollectionInstallOptionsFunc, 0)
			options := &collection.AnsibleGalaxyCollectionInstallOptions{}

			if len(parameters.Dependencies.Collections.Collections) > 0 {
				optionsFuncs = append(optionsFuncs, collection.WithCollectionNames(parameters.Dependencies.Collections.Collections...))
			}

			if len(parameters.Dependencies.Collections.APIKey) > 0 {
				options.APIKey = parameters.Dependencies.Collections.APIKey
			}

			options.ForceWithDeps = parameters.Dependencies.Collections.ForceWithDeps
			options.Pre = parameters.Dependencies.Collections.Pre

			if len(parameters.Dependencies.Collections.Timeout) > 0 {
				options.Timeout = parameters.Dependencies.Collections.Timeout
			}

			if len(parameters.Dependencies.Collections.Token) > 0 {
				options.Token = parameters.Dependencies.Collections.Token
			}

			options.IgnoreErrors = parameters.Dependencies.Collections.IgnoreErrors

			if len(parameters.Dependencies.Collections.RequirementsFile) > 0 {
				options.RequirementsFile = parameters.Dependencies.Collections.RequirementsFile
			}

			if len(parameters.Dependencies.Collections.Server) > 0 {
				options.Server = parameters.Dependencies.Collections.Server
			}

			options.Verbose = parameters.Dependencies.Collections.Verbose

			optionsFuncs = append(optionsFuncs, collection.WithGalaxyCollectionInstallOptions(options))
			galaxyInstallCollectionCmd := collection.NewAnsibleGalaxyCollectionInstallCmd(optionsFuncs...)

			galaxyInstallCollectionExecutor := configuration.NewAnsibleWithConfigurationSettingsExecute(
				execute.NewDefaultExecute(
					execute.WithCmd(galaxyInstallCollectionCmd),
					execute.WithCmdRunDir(workingDir),
				),
				configuration.WithAnsibleCollectionsPaths(filepath.Join(workingDir, CollectionsPath)),
			)

			workflowTasks = append(workflowTasks, galaxyInstallCollectionExecutor)
		}
	}

	playbookCmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybooks(parameters.Playbooks...),
		playbook.WithPlaybookOptions(ansiblePlaybookOptions),
	)

	playbookExecutor := configuration.NewAnsibleWithConfigurationSettingsExecute(
		execute.NewDefaultExecute(
			execute.WithCmd(playbookCmd),
			execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
			execute.WithCmdRunDir(workingDir),
		),
		configuration.WithAnsibleCollectionsPaths(filepath.Join(workingDir, CollectionsPath)),
	)

	workflowTasks = append(workflowTasks, playbookExecutor)

	workflowExecutor := workflow.NewWorkflowExecute(workflowTasks...)
	err := workflowExecutor.Execute(ctx)
	if err != nil {
		return err
	}

	return nil
}
