package executor

import (
	"path/filepath"
	"testing"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/execute/configuration"
	collection "github.com/apenella/go-ansible/v2/pkg/galaxy/collection/install"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestAnsiblePlaybookOptionsMapper(t *testing.T) {
	tests := []struct {
		desc string
		in   *entity.AnsiblePlaybookParameters
		out  *playbook.AnsiblePlaybookOptions
	}{
		{
			desc: "Testing AnsiblePlaybookOptionsMapper with all parameters",
			in: &entity.AnsiblePlaybookParameters{
				Check:             true,
				Diff:              true,
				ExtraVars:         map[string]interface{}{"key1": "value1", "key2": "value2"},
				ExtraVarsFile:     []string{"file1", "file2"},
				FlushCache:        true,
				ForceHandlers:     true,
				Forks:             10,
				Inventory:         "inventory",
				Limit:             "limit",
				ListHosts:         true,
				ListTags:          true,
				ListTasks:         true,
				SkipTags:          "skip",
				StartAtTask:       "task",
				SyntaxCheck:       true,
				Tags:              "tags",
				VaultID:           "vault",
				VaultPasswordFile: "vault-password-file",
				Verbose:           true,
				Version:           true,
				Connection:        "connection",
				SCPExtraArgs:      "scp-extra-args",
				SFTPExtraArgs:     "sftp-extra-args",
				SSHCommonArgs:     "ssh-common-args",
				SSHExtraArgs:      "ssh-extra-args",
				Timeout:           10,
				User:              "user",
				Become:            true,
				BecomeMethod:      "become-method",
				BecomeUser:        "become-user",
			},
			out: &playbook.AnsiblePlaybookOptions{
				Check:             true,
				Diff:              true,
				ExtraVars:         map[string]interface{}{"key1": "value1", "key2": "value2"},
				ExtraVarsFile:     []string{"file1", "file2"},
				FlushCache:        true,
				ForceHandlers:     true,
				Forks:             "10",
				Inventory:         "inventory",
				Limit:             "limit",
				ListHosts:         true,
				ListTags:          true,
				ListTasks:         true,
				SkipTags:          "skip",
				StartAtTask:       "task",
				SyntaxCheck:       true,
				Tags:              "tags",
				VaultID:           "vault",
				VaultPasswordFile: "vault-password-file",
				Verbose:           true,
				Version:           true,
				Connection:        "connection",
				SCPExtraArgs:      "scp-extra-args",
				SFTPExtraArgs:     "sftp-extra-args",
				SSHCommonArgs:     "ssh-common-args",
				SSHExtraArgs:      "ssh-extra-args",
				Timeout:           10,
				User:              "user",
				Become:            true,
				BecomeMethod:      "become-method",
				BecomeUser:        "become-user",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := ansiblePlaybookOptionsMapper(test.in)
			assert.Equal(t, test.out, res)
		})
	}

}

func TestAnsibleGalaxyCollectionInstallOptionsMapper(t *testing.T) {
	tests := []struct {
		desc string
		in   *entity.AnsiblePlaybookCollectionDependencies
		out  *collection.AnsibleGalaxyCollectionInstallOptions
	}{
		{
			desc: "Testing AnsibleGalaxyCollectionInstallOptionsMapper with all parameters",
			in: &entity.AnsiblePlaybookCollectionDependencies{
				APIKey:           "api-key",
				ForceWithDeps:    true,
				Pre:              true,
				Timeout:          "10",
				Token:            "token",
				IgnoreErrors:     true,
				RequirementsFile: "requirements-file",
				Server:           "server",
				Verbose:          true,
			},
			out: &collection.AnsibleGalaxyCollectionInstallOptions{
				APIKey:           "api-key",
				ForceWithDeps:    true,
				Pre:              true,
				Timeout:          "10",
				Token:            "token",
				IgnoreErrors:     true,
				RequirementsFile: "requirements-file",
				Server:           "server",
				Verbose:          true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := ansibleGalaxyCollectionInstallOptionsMapper(test.in)
			assert.Equal(t, test.out, res)
		})
	}

}

func TestCreateGalaxyCollectionInstallExecutor(t *testing.T) {
	run := NewAnsiblePlaybook(
		logger.NewFakeLogger(),
	)
	tests := []struct {
		desc       string
		workingDir string
		run        *AnsiblePlaybook
		in         *entity.AnsiblePlaybookParameters
		out        *configuration.AnsibleWithConfigurationSettingsExecute
	}{
		{
			desc:       "Testing creating a GalaxyCollectionInstallExecutor when parameters are not provided",
			run:        run,
			workingDir: "/tmp",
			in:         nil,
			out:        nil,
		},
		{
			desc:       "Testing creating a GalaxyCollectionInstallExecutor when working directory is not provided",
			run:        run,
			workingDir: "",
			in:         &entity.AnsiblePlaybookParameters{},
			out:        nil,
		},

		{
			desc:       "Testing creating a GalaxyCollectionInstallExecutor when dependency parameters are not provided",
			run:        run,
			workingDir: "/tmp",
			in:         &entity.AnsiblePlaybookParameters{},
			out:        nil,
		},
		{
			desc:       "Testing creating a GalaxyCollectionInstallExecutor when dependency parameters are provided and collections is not provided",
			run:        run,
			workingDir: "/tmp",
			in: &entity.AnsiblePlaybookParameters{
				Dependencies: &entity.AnsiblePlaybookDependencies{},
			},
			out: nil,
		},
		{
			desc:       "Testing creating a GalaxyCollectionInstallExecutor when collection names are provided",
			run:        run,
			workingDir: "/tmp",
			in: &entity.AnsiblePlaybookParameters{
				Dependencies: &entity.AnsiblePlaybookDependencies{
					Collections: &entity.AnsiblePlaybookCollectionDependencies{
						Collections: []string{"collection1", "collection2"},
					},
				},
			},
			out: configuration.NewAnsibleWithConfigurationSettingsExecute(
				execute.NewDefaultExecute(
					execute.WithCmd(
						collection.NewAnsibleGalaxyCollectionInstallCmd(
							[]collection.AnsibleGalaxyCollectionInstallOptionsFunc{
								collection.WithCollectionNames("collection1", "collection2"),
								collection.WithGalaxyCollectionInstallOptions(&collection.AnsibleGalaxyCollectionInstallOptions{}),
							}...,
						),
					),
					execute.WithCmdRunDir("/tmp"),
				),
				configuration.WithAnsibleCollectionsPaths(
					filepath.Join("/tmp", CollectionsPath),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.run.createGalaxyCollectionInstallExecutor(test.workingDir, test.in)
			assert.Equal(t, test.out, res)
		})
	}
}

func TestCreateAnsiblePlaybookExecutor(t *testing.T) {
	run := NewAnsiblePlaybook(
		logger.NewFakeLogger(),
	)

	tests := []struct {
		desc       string
		run        *AnsiblePlaybook
		workingDir string
		in         *entity.AnsiblePlaybookParameters
		out        *configuration.AnsibleWithConfigurationSettingsExecute
	}{
		{
			desc:       "Testing creating a AnsiblePlaybookExecutor when parameters are not provided",
			run:        run,
			workingDir: "/tmp",
			in:         nil,
			out:        nil,
		},
		{
			desc:       "Testing creating a AnsiblePlaybookExecutor when working directory is not provided",
			run:        run,
			workingDir: "",
			in:         &entity.AnsiblePlaybookParameters{},
			out:        nil,
		},
		{
			desc:       "Testing creating a AnsiblePlaybookExecutor when playbook list is not provided",
			run:        run,
			workingDir: "/tmp",
			in: &entity.AnsiblePlaybookParameters{
				Playbooks: nil,
			},
			out: nil,
		},
		{
			desc:       "Testing creating a AnsiblePlaybookExecutor when parameters are provided",
			run:        run,
			workingDir: "/tmp",
			in: &entity.AnsiblePlaybookParameters{
				Playbooks: []string{"playbook.yml"},
			},
			out: configuration.NewAnsibleWithConfigurationSettingsExecute(
				execute.NewDefaultExecute(
					execute.WithCmd(
						playbook.NewAnsiblePlaybookCmd(
							playbook.WithPlaybooks([]string{"playbook.yml"}...),
							playbook.WithPlaybookOptions(&playbook.AnsiblePlaybookOptions{}),
						),
					),
					execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
					execute.WithCmdRunDir("/tmp"),
				),
				configuration.WithAnsibleCollectionsPaths(
					filepath.Join("/tmp", CollectionsPath),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.run.createAnsiblePlaybookExecutor(test.workingDir, test.in)
			assert.Equal(t, test.out, res)
		})
	}
}
