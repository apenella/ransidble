package executor

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/execute/configuration"
	"github.com/apenella/go-ansible/v2/pkg/execute/workflow"
	collection "github.com/apenella/go-ansible/v2/pkg/galaxy/collection/install"
	role "github.com/apenella/go-ansible/v2/pkg/galaxy/role/install"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

const (
	// CollectionsPath represents the path where the collections are stored
	CollectionsPath = ".collections"
	// RolesPath represents the path where the roles are stored
	RolesPath = ".roles"
)

var (
	// ErrWorkingDirNotProvided represents an error when the working directory is not provided
	ErrWorkingDirNotProvided = fmt.Errorf("working directory not provided")
	// ErrParametersNotProvided represents an error when the parameters are not provided
	ErrParametersNotProvided = fmt.Errorf("parameters not provided")
	// ErrRunningAnsiblePlaybook represents an error when running an ansible playbook
	ErrRunningAnsiblePlaybook = fmt.Errorf("error running ansible playbook")
)

// AnsiblePlaybook represents an executor for running ansible playbooks
type AnsiblePlaybook struct {
	// logger is the logger
	logger repository.Logger
}

// NewAnsiblePlaybook returns a new AnsiblePlaybook instance
func NewAnsiblePlaybook(logger repository.Logger) *AnsiblePlaybook {
	return &AnsiblePlaybook{
		logger: logger,
	}
}

// Run runs an ansible playbook
func (a *AnsiblePlaybook) Run(ctx context.Context, workingDir string, parameters *entity.AnsiblePlaybookParameters) error {

	if workingDir == "" {
		a.logger.Error(
			ErrWorkingDirNotProvided.Error(),
			map[string]interface{}{
				"component": "AnsiblePlaybook.Run",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return ErrWorkingDirNotProvided
	}

	if parameters == nil {
		a.logger.Error(
			ErrParametersNotProvided.Error(),
			map[string]interface{}{
				"component": "AnsiblePlaybook.Run",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return ErrParametersNotProvided
	}

	workflowTasks := make([]execute.Executor, 0)
	playbookExecutor := a.createAnsiblePlaybookExecutor(workingDir, parameters)

	galaxyInstallCollectionExecutor := a.createGalaxyCollectionInstallExecutor(workingDir, parameters)
	if galaxyInstallCollectionExecutor != nil {
		workflowTasks = append(workflowTasks, galaxyInstallCollectionExecutor)
	}

	galaxyInstallRoleExecutor := a.createGalaxyRoleInstallExecutor(workingDir, parameters)
	if galaxyInstallRoleExecutor != nil {
		workflowTasks = append(workflowTasks, galaxyInstallRoleExecutor)
	}

	workflowTasks = append(workflowTasks, playbookExecutor)

	workflowExecutor := workflow.NewWorkflowExecute(workflowTasks...)
	err := workflowExecutor.Execute(ctx)
	if err != nil {
		a.logger.Error(
			fmt.Sprintf("%s: %s", ErrRunningAnsiblePlaybook, err),
			map[string]interface{}{
				"component": "AnsiblePlaybook.Run",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return fmt.Errorf("%s: %w", ErrRunningAnsiblePlaybook, err)
	}

	return nil
}

func (a *AnsiblePlaybook) createGalaxyRoleInstallExecutor(workingDir string, parameters *entity.AnsiblePlaybookParameters) *configuration.AnsibleWithConfigurationSettingsExecute {
	var galaxyInstallRolesExecutor *configuration.AnsibleWithConfigurationSettingsExecute

	if parameters == nil {
		a.logger.Debug(
			"Parameters not provided",
			map[string]interface{}{
				"component": "AnsiblePlaybook.createGalaxyRolesInstallExecutor",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return nil
	}

	if workingDir == "" {
		a.logger.Debug(
			"Working directory not provided",
			map[string]interface{}{
				"component": "AnsiblePlaybook.createGalaxyRolesInstallExecutor",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return nil
	}

	optionsFuncs := make([]role.AnsibleGalaxyRoleInstallOptionsFunc, 0)

	if parameters.Requirements != nil {
		if parameters.Requirements.Roles != nil {

			if len(parameters.Requirements.Roles.Roles) > 0 {
				optionsFuncs = append(optionsFuncs, role.WithRoleNames(parameters.Requirements.Roles.Roles...))
			}

			options := ansibleGalaxyRolesInstallOptionsMapper(parameters.Requirements.Roles)

			optionsFuncs = append(optionsFuncs, role.WithGalaxyRoleInstallOptions(options))
			galaxyInstallRolesCmd := role.NewAnsibleGalaxyRoleInstallCmd(optionsFuncs...)

			galaxyInstallRolesExecutor = configuration.NewAnsibleWithConfigurationSettingsExecute(
				execute.NewDefaultExecute(
					execute.WithCmd(galaxyInstallRolesCmd),
					execute.WithCmdRunDir(workingDir),
				),
				configuration.WithAnsibleRolesPath(filepath.Join(workingDir, RolesPath)),
			)
		}
	}

	return galaxyInstallRolesExecutor
}

func ansibleGalaxyRolesInstallOptionsMapper(parameters *entity.AnsiblePlaybookRoleRequirements) *role.AnsibleGalaxyRoleInstallOptions {

	options := &role.AnsibleGalaxyRoleInstallOptions{}

	if len(parameters.APIKey) > 0 {
		options.ApiKey = parameters.APIKey
	}

	options.IgnoreErrors = parameters.IgnoreErrors
	options.NoDeps = parameters.NoDeps

	if len(parameters.RoleFile) > 0 {
		options.RoleFile = parameters.RoleFile
	}

	if len(parameters.Server) > 0 {
		options.Server = parameters.Server
	}

	if len(parameters.Timeout) > 0 {
		options.Timeout = parameters.Timeout
	}

	if len(parameters.Token) > 0 {
		options.Token = parameters.Token
	}

	options.Verbose = parameters.Verbose

	return options
}

// createGalaxyCollectionInstallExecutor returns an Executor to run the Ansible Galaxy Collection install command
func (a *AnsiblePlaybook) createGalaxyCollectionInstallExecutor(workingDir string, parameters *entity.AnsiblePlaybookParameters) *configuration.AnsibleWithConfigurationSettingsExecute {

	var galaxyInstallCollectionExecutor *configuration.AnsibleWithConfigurationSettingsExecute

	if parameters == nil {
		a.logger.Debug(
			"Parameters not provided",
			map[string]interface{}{
				"component": "AnsiblePlaybook.createGalaxyCollectionInstallExecutor",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return nil
	}

	if workingDir == "" {
		a.logger.Debug(
			"Working directory not provided",
			map[string]interface{}{
				"component": "AnsiblePlaybook.createGalaxyCollectionInstallExecutor",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return nil
	}

	optionsFuncs := make([]collection.AnsibleGalaxyCollectionInstallOptionsFunc, 0)

	if parameters.Requirements != nil {
		if parameters.Requirements.Collections != nil {

			if len(parameters.Requirements.Collections.Collections) > 0 {
				optionsFuncs = append(optionsFuncs, collection.WithCollectionNames(parameters.Requirements.Collections.Collections...))
			}

			options := ansibleGalaxyCollectionInstallOptionsMapper(parameters.Requirements.Collections)

			optionsFuncs = append(optionsFuncs, collection.WithGalaxyCollectionInstallOptions(options))
			galaxyInstallCollectionCmd := collection.NewAnsibleGalaxyCollectionInstallCmd(optionsFuncs...)

			galaxyInstallCollectionExecutor = configuration.NewAnsibleWithConfigurationSettingsExecute(
				execute.NewDefaultExecute(
					execute.WithCmd(galaxyInstallCollectionCmd),
					execute.WithCmdRunDir(workingDir),
				),
				configuration.WithAnsibleCollectionsPaths(filepath.Join(workingDir, CollectionsPath)),
			)
		}
	}
	return galaxyInstallCollectionExecutor
}

// createAnsiblePlaybookExecutor returns an Executor to run the Ansible Playbook command
func (a *AnsiblePlaybook) createAnsiblePlaybookExecutor(workingDir string, parameters *entity.AnsiblePlaybookParameters) *configuration.AnsibleWithConfigurationSettingsExecute {

	var playbookExecutor *configuration.AnsibleWithConfigurationSettingsExecute

	if parameters == nil {
		a.logger.Debug(
			"Parameters not provided",
			map[string]interface{}{
				"component": "AnsiblePlaybook.createAnsiblePlaybookExecutor",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return nil
	}

	if workingDir == "" {
		a.logger.Debug(
			"Working directory not provided",
			map[string]interface{}{
				"component": "AnsiblePlaybook.createAnsiblePlaybookExecutor",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			})

		return nil
	}

	// it is considered that if there are no playbooks to run, the executor is not created. However, the an error should be returned by the caller before calling this function
	if len(parameters.Playbooks) == 0 {
		return nil
	}

	ansiblePlaybookOptions := ansiblePlaybookOptionsMapper(parameters)

	playbookCmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybooks(parameters.Playbooks...),
		playbook.WithPlaybookOptions(ansiblePlaybookOptions),
	)

	playbookExecutor = configuration.NewAnsibleWithConfigurationSettingsExecute(
		execute.NewDefaultExecute(
			execute.WithCmd(playbookCmd),
			execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
			execute.WithCmdRunDir(workingDir),
		),
		configuration.WithAnsibleCollectionsPaths(filepath.Join(workingDir, CollectionsPath)),
	)

	return playbookExecutor
}

// ansiblePlaybookOptionsMapper maps an entity.AnsiblePlaybookParameters to a playbook.AnsiblePlaybookOptions
func ansiblePlaybookOptionsMapper(parameters *entity.AnsiblePlaybookParameters) *playbook.AnsiblePlaybookOptions {

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{}

	ansiblePlaybookOptions.Check = parameters.Check

	ansiblePlaybookOptions.Diff = parameters.Diff

	if len(parameters.ExtraVars) > 0 {
		if ansiblePlaybookOptions.ExtraVars == nil {
			ansiblePlaybookOptions.ExtraVars = make(map[string]interface{})
		}

		for k, v := range parameters.ExtraVars {
			ansiblePlaybookOptions.ExtraVars[k] = v
		}
	}

	if len(parameters.ExtraVarsFile) > 0 {
		ansiblePlaybookOptions.ExtraVarsFile = append([]string{}, parameters.ExtraVarsFile...)
	}

	ansiblePlaybookOptions.FlushCache = parameters.FlushCache

	ansiblePlaybookOptions.ForceHandlers = parameters.ForceHandlers

	if parameters.Forks > 0 {
		ansiblePlaybookOptions.Forks = strconv.Itoa(parameters.Forks)
	}

	if len(parameters.Inventory) > 0 {
		ansiblePlaybookOptions.Inventory = parameters.Inventory
	}

	if len(parameters.Limit) > 0 {
		ansiblePlaybookOptions.Limit = parameters.Limit
	}

	ansiblePlaybookOptions.ListHosts = parameters.ListHosts
	ansiblePlaybookOptions.ListTags = parameters.ListTags
	ansiblePlaybookOptions.ListTasks = parameters.ListTasks

	if len(parameters.SkipTags) > 0 {
		ansiblePlaybookOptions.SkipTags = parameters.SkipTags
	}

	if len(parameters.StartAtTask) > 0 {
		ansiblePlaybookOptions.StartAtTask = parameters.StartAtTask
	}

	ansiblePlaybookOptions.SyntaxCheck = parameters.SyntaxCheck

	if len(parameters.Tags) > 0 {
		ansiblePlaybookOptions.Tags = parameters.Tags
	}

	if len(parameters.VaultID) > 0 {
		ansiblePlaybookOptions.VaultID = parameters.VaultID
	}

	if len(parameters.VaultPasswordFile) > 0 {
		ansiblePlaybookOptions.VaultPasswordFile = parameters.VaultPasswordFile
	}

	ansiblePlaybookOptions.Verbose = parameters.Verbose

	ansiblePlaybookOptions.Version = parameters.Version

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

	ansiblePlaybookOptions.Become = parameters.Become

	if len(parameters.BecomeMethod) > 0 {
		ansiblePlaybookOptions.BecomeMethod = parameters.BecomeMethod
	}

	if len(parameters.BecomeUser) > 0 {
		ansiblePlaybookOptions.BecomeUser = parameters.BecomeUser
	}

	return ansiblePlaybookOptions
}

func ansibleGalaxyCollectionInstallOptionsMapper(parameters *entity.AnsiblePlaybookCollectionRequirements) *collection.AnsibleGalaxyCollectionInstallOptions {

	options := &collection.AnsibleGalaxyCollectionInstallOptions{}

	if len(parameters.APIKey) > 0 {
		options.APIKey = parameters.APIKey
	}

	options.ForceWithDeps = parameters.ForceWithDeps
	options.Pre = parameters.Pre

	if len(parameters.Timeout) > 0 {
		options.Timeout = parameters.Timeout
	}

	if len(parameters.Token) > 0 {
		options.Token = parameters.Token
	}

	options.IgnoreErrors = parameters.IgnoreErrors

	if len(parameters.RequirementsFile) > 0 {
		options.RequirementsFile = parameters.RequirementsFile
	}

	if len(parameters.Server) > 0 {
		options.Server = parameters.Server
	}

	options.Verbose = parameters.Verbose

	return options
}
