package functional

import (
	"context"
	"fmt"
	"io"
	nethttp "net/http"
	"testing"
	"time"

	projectService "github.com/apenella/ransidble/internal/domain/core/service/project"
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
	listenAddress    string
	openAPIValidator *OpenAPIValidator
	router           *echo.Echo
	server           *http.Server

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *SuiteGetProject) SetupSuite() {
	var err error

	suite.openAPIValidator, err = PrepareOpenAPIValidator(openAPIDefPath)
	if err != nil {
		suite.T().Errorf("Error initializing OpenAPI validator: %s", err)
		suite.T().FailNow()
		return
	}
}

// SetupTest runs before each test in the suite
func (suite *SuiteGetProject) SetupTest() {
	suite.listenAddress = "0.0.0.0:8080"
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
func (suite *SuiteGetProject) TestGetProjectProject1() {

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

	if suite.openAPIValidator == nil {
		suite.T().Errorf("%s. OpenAPI validator is not initialized", suite.T().Name())
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
		expectedStatusCode int
		arrangeTest        func(*SuiteGetProject)
	}{
		{
			desc:   "Functional test to get an existing project",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-1",
			arrangeTest: func(suite *SuiteGetProject) {
				log := logger.NewFakeLogger()
				afs := afero.NewOsFs()

				projectsRepository := localprojectpersistence.NewLocalProjectRepository(
					afs,
					"../projects",
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
			expectedStatusCode: nethttp.StatusOK,
		},
		{
			desc:   "Functional test to get a non-existing project",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-non-existing",
			arrangeTest: func(suite *SuiteGetProject) {
				log := logger.NewFakeLogger()
				afs := afero.NewOsFs()

				projectsRepository := localprojectpersistence.NewLocalProjectRepository(
					afs,
					"../projects",
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
			expectedStatusCode: nethttp.StatusNotFound,
		},
	}

	for _, test := range tests {

		var err error
		var httpReq *nethttp.Request
		var httpResp *nethttp.Response
		var body []byte

		httpReq, err = nethttp.NewRequest(test.method, test.url, nil)
		if err != nil {
			suite.T().Errorf("%s. Error creating HTTP request: %s", suite.T().Name(), err)
			suite.T().FailNow()
			return
		}

		suite.T().Run(fmt.Sprintf("functional %s", test.desc), func(t *testing.T) {
			suite.T().Log("Functional: " + test.desc)

			if test.arrangeTest != nil {
				test.arrangeTest(suite)
			}

			client := &nethttp.Client{}
			httpResp, err = client.Do(httpReq)
			if err != nil {
				suite.T().Errorf("%s. Error performing HTTP request: %s", suite.T().Name(), err)
				suite.T().FailNow()
				return
			}
			defer httpResp.Body.Close()

			body, err = io.ReadAll(httpResp.Body)
			if err != nil {
				suite.T().Errorf("%s. Error reading response body: %s", suite.T().Name(), err)
				suite.T().FailNow()
				return
			}

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), test.expectedStatusCode, httpResp.StatusCode)
		})

		suite.T().Run(fmt.Sprintf("openapi %s", test.desc), func(t *testing.T) {
			suite.T().Log("OpenAPI: " + test.desc)
			err = suite.openAPIValidator.ValidateResponse(body, httpReq, httpResp.StatusCode, httpResp.Header)
			assert.NoError(suite.T(), err)
		})
	}
}

// TestFunctionalGetProject runs the test suite
func TestFunctionalGetProject(t *testing.T) {
	suite.Run(t, new(SuiteGetProject))
}
