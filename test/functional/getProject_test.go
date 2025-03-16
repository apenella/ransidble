package functional

import (
	"context"
	"errors"
	nethttp "net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	projectService "github.com/apenella/ransidble/internal/domain/core/service/project"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	serve "github.com/apenella/ransidble/internal/handler/cli/serve"
	"github.com/apenella/ransidble/internal/handler/http"
	projectHandler "github.com/apenella/ransidble/internal/handler/http/project"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	localprojectpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository"
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SuiteGetProject is the test suite for the HTTP server
type SuiteGetProject struct {
	listenAddress string
	router        *echo.Echo
	server        *http.Server

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *SuiteGetProject) SetupSuite() {
	suite.listenAddress = "0.0.0.0:8080"
}

// SetupTest runs before each test in the suite
func (suite *SuiteGetProject) SetupTest() {
	suite.router = echo.New()
	suite.server = http.NewServer(suite.listenAddress, suite.router, logger.NewFakeLogger())
}

// TearDownSuite runs after all tests in this suite have run
func (suite *SuiteGetProject) TearDownSuite() {
	suite.server.Stop()
}

// TearDownTest runs after each test in the suite
func (suite *SuiteGetProject) TearDownTest() {}

// TestGetProjectProject1 is a functional test for the GetProject endpoint
func (suite *SuiteGetProject) TestGetProject() {

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
	defer suite.server.Stop()

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
		expectedBody       string
		expectedStatusCode int
		arrangeTest        func(*SuiteGetProject)
	}{
		{
			desc:   "Testing a request to get project functinal behavior when project exists that returns a StatusOK status code",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-1",
			arrangeTest: func(suite *SuiteGetProject) {
				log := logger.NewFakeLogger()
				afs := afero.NewOsFs()

				projectsRepository := localprojectpersistence.NewLocalProjectRepository(
					afs,
					filepath.Join("..", "fixtures", "functional-get-project"),
					log,
				)

				errLoadProjects := projectsRepository.LoadProjects()
				if errLoadProjects != nil {
					suite.T().Errorf("Error loading projects: %s", errLoadProjects)
					suite.T().FailNow()
					return
				}

				getProjectService := projectService.NewGetProjectService(projectsRepository, log)
				getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
				suite.router.GET(serve.GetProjectPath, getProjectHandler.Handle)
			},
			expectedBody:       "{\"format\":\"plain\",\"name\":\"project-1\",\"reference\":\"../fixtures/functional-get-project/project-1\",\"storage\":\"local\"}",
			expectedStatusCode: nethttp.StatusOK,
		},
		{
			desc:   "Testing a request to get project functional behavior when project does not exist that returns a StatusNotFound status code",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-non-existing",
			arrangeTest: func(suite *SuiteGetProject) {
				log := logger.NewFakeLogger()
				afs := afero.NewOsFs()

				projectsRepository := localprojectpersistence.NewLocalProjectRepository(
					afs,
					filepath.Join("..", "fixtures", "functional-get-project"),
					log,
				)

				errLoadProjects := projectsRepository.LoadProjects()
				if errLoadProjects != nil {
					suite.T().Errorf("Error loading projects: %s", errLoadProjects)
					suite.T().FailNow()
					return
				}

				getProjectService := projectService.NewGetProjectService(projectsRepository, log)
				getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
				suite.router.GET(serve.GetProjectPath, getProjectHandler.Handle)
			},
			expectedBody:       "{\"id\":\"\",\"error\":\"error getting project: error finding project: project not found\"}",
			expectedStatusCode: nethttp.StatusNotFound,
		},
		{
			desc:   "Testing a request get project functional behaviour when there is an internal error that returns a StatusInternalServerError status code",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-internal-error",
			arrangeTest: func(suite *SuiteGetProject) {
				log := logger.NewFakeLogger()
				afs := afero.NewOsFs()

				projectsRepository := localprojectpersistence.NewLocalProjectRepository(
					afs,
					filepath.Join("..", "fixtures", "functional-get-project"),
					log,
				)

				errLoadProjects := projectsRepository.LoadProjects()
				if errLoadProjects != nil {
					suite.T().Errorf("Error loading projects: %s", errLoadProjects)
					suite.T().FailNow()
					return
				}

				// This is the mock service that simulates an internal error returned by the service
				getProjectService := service.NewMockGetProjectService()
				getProjectService.On("GetProject", "project-internal-error").Return(
					&entity.Project{},
					errors.New("testing get project internal error"),
				)

				getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
				suite.router.GET(serve.GetProjectPath, getProjectHandler.Handle)
			},
			expectedBody:       "{\"id\":\"\",\"error\":\"error getting project: testing get project internal error\"}",
			expectedStatusCode: nethttp.StatusInternalServerError,
		},
		{
			desc:   "Testing a request get project functional behaviour when service returns a project not provided error that returns a StatusBadRequest status code",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-not-provided",
			arrangeTest: func(suite *SuiteGetProject) {
				log := logger.NewFakeLogger()
				afs := afero.NewOsFs()

				projectsRepository := localprojectpersistence.NewLocalProjectRepository(
					afs,
					filepath.Join("..", "fixtures", "functional-get-project"),
					log,
				)

				errLoadProjects := projectsRepository.LoadProjects()
				if errLoadProjects != nil {
					suite.T().Errorf("Error loading projects: %s", errLoadProjects)
					suite.T().FailNow()
					return
				}

				// This is the mock service that simulates a project not provided error returned by the service
				getProjectService := service.NewMockGetProjectService()
				getProjectService.On("GetProject", "project-not-provided").Return(
					&entity.Project{},
					domainerror.NewProjectNotProvidedError(errors.New("testing get project not provided error")),
				)

				getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
				suite.router.GET(serve.GetProjectPath, getProjectHandler.Handle)
			},
			expectedBody:       "{\"id\":\"\",\"error\":\"error getting project: testing get project not provided error\"}",
			expectedStatusCode: nethttp.StatusBadRequest,
		},
	}

	for _, test := range tests {

		if test.arrangeTest != nil {
			test.arrangeTest(suite)
		}

		input := &InputFunctionalTest{
			desc:               test.desc,
			method:             test.method,
			url:                test.url,
			expectedStatusCode: test.expectedStatusCode,
			expectedBody:       test.expectedBody,
		}

		err := actAndAssert(suite.T(), input)
		assert.NoError(suite.T(), err)
	}
}

// TestFunctionalGetProject runs the test suite
func TestFunctionalGetProject(t *testing.T) {
	suite.Run(t, new(SuiteGetProject))
}
