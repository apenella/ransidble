package functional

import (
	"context"
	"fmt"
	"io"
	nethttp "net/http"
	"strings"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/service/executor"
	taskService "github.com/apenella/ransidble/internal/domain/core/service/task"
	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
	"github.com/apenella/ransidble/internal/handler/cli/serve"
	"github.com/apenella/ransidble/internal/handler/http"
	taskHandler "github.com/apenella/ransidble/internal/handler/http/task"
	"github.com/apenella/ransidble/internal/infrastructure/filesystem"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch"
	localprojectpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository"
	taskpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
	"github.com/apenella/ransidble/internal/infrastructure/tar"
	"github.com/apenella/ransidble/internal/infrastructure/unpack"
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// SuiteCreateTaskAnsiblePlaybook is the test suite for the HTTP server
type SuiteCreateTaskAnsiblePlaybook struct {
	listenAddress string
	router        *echo.Echo
	server        *http.Server

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *SuiteCreateTaskAnsiblePlaybook) SetupSuite() {
	suite.listenAddress = "0.0.0.0:8080"
}

// SetupTest runs before each test in the suite
func (suite *SuiteCreateTaskAnsiblePlaybook) SetupTest() {
	suite.router = echo.New()
	suite.server = http.NewServer(suite.listenAddress, suite.router, logger.NewFakeLogger())
}

// TearDownSuite runs after all tests in this suite have run
func (suite *SuiteCreateTaskAnsiblePlaybook) TearDownSuite() {
	suite.server.Stop()
}

// TestTaskAnsiblePlaybook is a functional test for the TaskAnsiblePlaybook endpoint
func (suite *SuiteCreateTaskAnsiblePlaybook) TestCreateTaskAnsiblePlaybook() {

	if suite.server == nil {
		suite.T().Errorf("%s. HTTP server is not initialized", suite.T().Name())
		suite.T().FailNow()
		return
	}

	if suite.router == nil {
		suite.T().Errorf("%s. HTTP router is not initialized", suite.T().Name())
		suite.T().FailNow()
		return
	}

	if suite.listenAddress == "" {
		suite.T().Errorf("%s. Listen address is not initialized", suite.T().Name())
		suite.T().FailNow()
		return
	}

	go func() {
		err := suite.server.Start(context.Background())
		if err != nil {
			suite.T().Errorf("%s. error starting HTTP server: %s", suite.T().Name(), err)
			suite.T().FailNow()
			return
		}
	}()

	errConn := waitHTTPServer(suite.listenAddress, 1*time.Second, 5)
	if errConn != nil {
		suite.T().Errorf("%s. error waiting for HTTP server: %s", suite.T().Name(), errConn)
		suite.T().FailNow()
		return
	}

	tests := []struct {
		desc               string
		method             string
		url                string
		expectedStatusCode int
		// the function returns the dispatcher as a workaround to start the dispatcher within the test function
		arrangeTest func(*SuiteCreateTaskAnsiblePlaybook) (*executor.Dispatch, error)
		parameters  io.ReadCloser
	}{
		{
			desc:       "Testing a request to create an Ansible Playbook task successfully that returns a StatusAccepted status code",
			method:     "POST",
			url:        "http://" + suite.listenAddress + "/tasks/ansible-playbook/project-1",
			parameters: io.NopCloser(strings.NewReader(`{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}`)),
			arrangeTest: func(suite *SuiteCreateTaskAnsiblePlaybook) (*executor.Dispatch, error) {
				// ansibleExecutorMock is a mock dependency to avoid executing the playbook on every test
				ansibleExecutor := executor.NewMockAnsiblePlaybookExecutor()

				expectedParameters := &entity.AnsiblePlaybookParameters{
					Playbooks:         []string{"site.yml"},
					Check:             false,
					Diff:              false,
					Requirements:      &entity.AnsiblePlaybookRequirements{},
					ExtraVars:         map[string]interface{}{},
					ExtraVarsFile:     []string{},
					FlushCache:        false,
					ForceHandlers:     false,
					Forks:             0,
					Inventory:         "127.0.0.1,",
					Limit:             "",
					ListHosts:         false,
					ListTags:          false,
					ListTasks:         false,
					SkipTags:          "",
					StartAtTask:       "",
					SyntaxCheck:       false,
					Tags:              "",
					VaultID:           "",
					VaultPasswordFile: "",
					Verbose:           false,
					Version:           false,
					Connection:        "local",
					SCPExtraArgs:      "",
					SFTPExtraArgs:     "",
					SSHCommonArgs:     "",
					SSHExtraArgs:      "",
					Timeout:           0,
					User:              "",
					Become:            false,
					BecomeMethod:      "",
					BecomeUser:        "",
				}

				ansibleExecutor.On("Run", mock.Anything, mock.Anything, expectedParameters).Return(nil)

				dispatcher, err := arrangeTaskAnsiblePlaybookRouter(suite.router, ansibleExecutor)
				if err != nil {
					return nil, fmt.Errorf("error arranging router: %s", err)
				}

				return dispatcher, nil
			},
			expectedStatusCode: nethttp.StatusAccepted,
		},
		{
			desc:   "Testing a request to create an Ansible Playbook task successfully passing all the parameters that returns a StatusAccepted status code",
			method: "POST",
			url:    "http://" + suite.listenAddress + "/tasks/ansible-playbook/project-1",
			parameters: io.NopCloser(strings.NewReader(`
{
  "playbooks": ["playbook1.yml", "playbook2.yml"],
  "check": true,
  "diff": true,
  "requirements": {
    "roles": {
      "roles": ["role1", "role2"],
      "api_key": "your_api_key",
      "ignore_errors": true,
      "no_deps": true,
      "role_file": "roles/requirements.yml",
      "server": "https://galaxy.ansible.com",
      "timeout": 60,
      "token": "your_token",
      "verbose": true
    },
    "collections": {
      "collections": ["collection1", "collection2"],
      "api_key": "your_api_key",
      "force_with_deps": true,
      "pre": true,
      "timeout": 70,
      "token": "your_token",
      "ignore_errors": true,
      "requirements_file": "collections/requirements.yml",
      "server": "https://galaxy.ansible.com",
      "verbose": true
    }
  },
  "extra_vars": {
    "var1": "value1",
    "var2": "value2"
  },
  "extra_vars_file": ["extra_vars1.yml", "extra_vars2.yml"],
  "flush_cache": true,
  "force_handlers": true,
  "forks": 10,
  "inventory": "inventory/hosts",
  "limit": "all",
  "list_hosts": true,
  "list_tags": true,
  "list_tasks": true,
  "skip_tags": "tag1,tag2",
  "start_at_task": "task1",
  "syntax_check": true,
  "tags": "tag1,tag2",
  "vault_id": "vault_id",
  "vault_password_file": "vault_password_file",
  "verbose": true,
  "version": true,
  "connection": "ssh",
  "scp_extra_args": "scp_extra_args",
  "sftp_extra_args": "sftp_extra_args",
  "ssh_common_args": "ssh_common_args",
  "ssh_extra_args": "ssh_extra_args",
  "timeout": 30,
  "user": "ansible",
  "become": true,
  "become_method": "sudo",
  "become_user": "root"
}
`)),
			arrangeTest: func(suite *SuiteCreateTaskAnsiblePlaybook) (*executor.Dispatch, error) {
				// ansibleExecutorMock is a mock dependency to avoid executing the playbook on every test
				ansibleExecutor := executor.NewMockAnsiblePlaybookExecutor()

				expectedParameters := &entity.AnsiblePlaybookParameters{
					Playbooks: []string{"playbook1.yml", "playbook2.yml"},
					Check:     true,
					Diff:      true,
					Requirements: &entity.AnsiblePlaybookRequirements{
						Roles: &entity.AnsiblePlaybookRoleRequirements{
							Roles:        []string{"role1", "role2"},
							APIKey:       "your_api_key",
							IgnoreErrors: true,
							NoDeps:       true,
							RoleFile:     "roles/requirements.yml",
							Server:       "https://galaxy.ansible.com",
							Timeout:      60,
							Token:        "your_token",
							Verbose:      true,
						},
						Collections: &entity.AnsiblePlaybookCollectionRequirements{
							Collections:      []string{"collection1", "collection2"},
							APIKey:           "your_api_key",
							ForceWithDeps:    true,
							Pre:              true,
							Timeout:          70,
							Token:            "your_token",
							IgnoreErrors:     true,
							RequirementsFile: "collections/requirements.yml",
							Server:           "https://galaxy.ansible.com",
							Verbose:          true,
						},
					},
					ExtraVars: map[string]interface{}{
						"var1": "value1",
						"var2": "value2",
					},
					ExtraVarsFile:     []string{"extra_vars1.yml", "extra_vars2.yml"},
					FlushCache:        true,
					ForceHandlers:     true,
					Forks:             10,
					Inventory:         "inventory/hosts",
					Limit:             "all",
					ListHosts:         true,
					ListTags:          true,
					ListTasks:         true,
					SkipTags:          "tag1,tag2",
					StartAtTask:       "task1",
					SyntaxCheck:       true,
					Tags:              "tag1,tag2",
					VaultID:           "vault_id",
					VaultPasswordFile: "vault_password_file",
					Verbose:           true,
					Version:           true,
					Connection:        "ssh",
					SCPExtraArgs:      "scp_extra_args",
					SFTPExtraArgs:     "sftp_extra_args",
					SSHCommonArgs:     "ssh_common_args",
					SSHExtraArgs:      "ssh_extra_args",
					Timeout:           30,
					User:              "ansible",
					Become:            true,
					BecomeMethod:      "sudo",
					BecomeUser:        "root",
				}

				ansibleExecutor.On("Run", mock.Anything, mock.Anything, expectedParameters).Return(nil)

				dispatcher, err := arrangeTaskAnsiblePlaybookRouter(suite.router, ansibleExecutor)
				if err != nil {
					return nil, fmt.Errorf("error arranging router: %s", err)
				}

				return dispatcher, nil
			},
			expectedStatusCode: nethttp.StatusAccepted,
		},
		{
			desc:       "Testing a request to create an Ansible Playbook task for a non existing project that returns a StatusNotFound status code",
			method:     "POST",
			url:        "http://" + suite.listenAddress + "/tasks/ansible-playbook/non-existing-project",
			parameters: io.NopCloser(strings.NewReader(`{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}`)),
			arrangeTest: func(suite *SuiteCreateTaskAnsiblePlaybook) (*executor.Dispatch, error) {
				// ansibleExecutorMock is a mock dependency to avoid executing the playbook on every test
				ansibleExecutor := executor.NewMockAnsiblePlaybookExecutor()

				expectedParameters := &entity.AnsiblePlaybookParameters{
					Playbooks: []string{"site.yml"},

					Requirements:  &entity.AnsiblePlaybookRequirements{},
					ExtraVars:     map[string]interface{}{},
					ExtraVarsFile: []string{},
					Inventory:     "127.0.0.1,",
					Connection:    "local",
				}

				ansibleExecutor.On("Run", mock.Anything, mock.Anything, expectedParameters).Return(nil)

				dispatcher, err := arrangeTaskAnsiblePlaybookRouter(suite.router, ansibleExecutor)
				if err != nil {
					return nil, fmt.Errorf("error arranging router: %s", err)
				}

				return dispatcher, nil
			},
			expectedStatusCode: nethttp.StatusNotFound,
		},
		{
			desc:       "Testing a request to create an Ansible Playbook task with an invalid payload (missing inventory) that returns a StatusBadRequest status code",
			method:     "POST",
			url:        "http://" + suite.listenAddress + "/tasks/ansible-playbook/project-1",
			parameters: io.NopCloser(strings.NewReader(`{"playbooks": ["site.yml"], "connection": "local"}`)),
			arrangeTest: func(suite *SuiteCreateTaskAnsiblePlaybook) (*executor.Dispatch, error) {
				// ansibleExecutor is not required for this test. An error is expected before the executor is called
				dispatcher, err := arrangeTaskAnsiblePlaybookRouter(suite.router, nil)
				if err != nil {
					return nil, fmt.Errorf("error arranging router: %s", err)
				}

				return dispatcher, nil
			},
			expectedStatusCode: nethttp.StatusBadRequest,
		},
	}

	for _, test := range tests {

		if test.arrangeTest != nil {
			dispatcher, err := test.arrangeTest(suite)
			assert.NoError(suite.T(), err)

			go func() {
				err := dispatcher.Start(context.TODO())
				assert.NoError(suite.T(), err)
			}()
		}

		input := &InputFunctionalTest{
			desc:       test.desc,
			method:     test.method,
			url:        test.url,
			parameters: test.parameters,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedStatusCode: test.expectedStatusCode,
		}

		err := actAndAssert(suite.T(), input)
		assert.NoError(suite.T(), err)
	}
}

// TestFunctionalSuiteCreateTaskAnsiblePlaybook runs the test suite
func TestFunctionalSuiteCreateTaskAnsiblePlaybook(t *testing.T) {
	suite.Run(t, new(SuiteCreateTaskAnsiblePlaybook))
}

// arrangeTaskAnsiblePlaybookRouter arranges the router for the test. The function returns the dispatcher as a workaround to start the dispatcher within the test function
func arrangeTaskAnsiblePlaybookRouter(router *echo.Echo, ansibleExecutor executor.AnsiblePlaybookExecutor) (*executor.Dispatch, error) {
	log := logger.NewFakeLogger()

	roFsBase := afero.NewReadOnlyFs(afero.NewOsFs())
	rwFs := afero.NewCopyOnWriteFs(roFsBase, afero.NewMemMapFs())
	fs := filesystem.NewFilesystem(rwFs)

	projectsRepository := localprojectpersistence.NewLocalProjectRepository(
		rwFs,
		"../projects",
		log,
	)

	errLoadProjects := projectsRepository.LoadProjects()
	if errLoadProjects != nil {
		return nil, fmt.Errorf("error loading projects: %s", errLoadProjects)
	}

	fetchFactory := fetch.NewFactory()
	fetchFactory.Register(
		entity.ProjectTypeLocal,
		fetch.NewLocalStorage(
			rwFs,
			log,
		),
	)

	unpackFactory := unpack.NewFactory()
	unpackFactory.Register(entity.ProjectFormatPlain, unpack.NewPlainFormat(
		rwFs,
		log,
	))

	tarExtractor := tar.NewTar(rwFs, log)
	unpackFactory.Register(entity.ProjectFormatTarGz, unpack.NewTarGzipFormat(
		rwFs,
		tarExtractor,
		log,
	))

	workspaceBuilder := workspace.NewBuilder(
		fs,
		fetchFactory,
		unpackFactory,
		projectsRepository,
		log,
	)

	dispatcher := executor.NewDispatch(
		1,
		workspaceBuilder,
		ansibleExecutor,
		log,
	)

	taskRepository := taskpersistence.NewMemoryTaskRepository(log)
	createTaskAnsiblePlaybookService := taskService.NewCreateTaskAnsiblePlaybookService(
		dispatcher,
		taskRepository,
		projectsRepository,
		log,
	)
	createTaskAnsiblePlaybookHandler := taskHandler.NewCreateTaskAnsiblePlaybookHandler(createTaskAnsiblePlaybookService, log)

	router.POST(serve.CreateTaskAnsiblePlaybookPath, createTaskAnsiblePlaybookHandler.Handle)

	return dispatcher, nil
}
