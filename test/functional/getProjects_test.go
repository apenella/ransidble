package functional

import (
	"context"
	"errors"
	"fmt"
	"io"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
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

// SuiteGetProjects is the test suite for the HTTP server
type SuiteGetProjects struct {
	listenAddress    string
	openAPIValidator *OpenAPIValidator
	router           *echo.Echo
	server           *http.Server

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *SuiteGetProjects) SetupSuite() {
	var err error

	suite.openAPIValidator, err = PrepareOpenAPIValidator(openAPIDefPath)
	if err != nil {
		suite.T().Errorf("Error initializing OpenAPI validator: %s", err)
		suite.T().FailNow()
		return
	}
}

// TearDownSuite runs after all tests in this suite have run
func (suite *SuiteGetProjects) TearDownSuite() {}

// SetupTest runs before each test in the suite
func (suite *SuiteGetProjects) SetupTest() {
	suite.listenAddress = "0.0.0.0:8080"
	suite.router = echo.New()
	suite.server = http.NewServer(suite.listenAddress, suite.router, logger.NewFakeLogger())
}

// TearDownTest runs after each test in the suite
func (suite *SuiteGetProjects) TearDownTest() {
	suite.server.Stop()
}

// TestGetProjects is a functional test for the GetProjects endpoint
func (suite *SuiteGetProjects) TestGetProjects() {

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
		arrangeTest        func(*SuiteGetProjects)
	}{
		{
			desc:   "Testing get projects list functinoal behavior when the request is successful",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath,
			arrangeTest: func(suite *SuiteGetProjects) {
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
				getProjectListHandler := projectHandler.NewGetProjectListHandler(getProjectService, log)
				suite.router.GET(serve.GetProjectsPath, getProjectListHandler.Handle)
			},
			expectedStatusCode: nethttp.StatusOK,
		},
		{
			desc:   "Testing get projects list functional behaviour when service returns an error getting the projects list",
			method: nethttp.MethodGet,
			url:    "http://" + suite.listenAddress + serve.GetProjectsPath,
			arrangeTest: func(suite *SuiteGetProjects) {
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

				// This is the mock service that simulates a project not provided error returned by the service
				getProjectService := service.NewMockGetProjectService()
				getProjectService.On("GetProjectsList").Return(
					[]*entity.Project{},
					errors.New("testing get project list error"),
				)

				getProjectListHandler := projectHandler.NewGetProjectListHandler(getProjectService, log)
				suite.router.GET(serve.GetProjectsPath, getProjectListHandler.Handle)
			},
			expectedStatusCode: nethttp.StatusInternalServerError,
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

		suite.T().Run(fmt.Sprintf("Functional %s", test.desc), func(t *testing.T) {

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

		suite.T().Run(fmt.Sprintf("OpenAPI %s", test.desc), func(t *testing.T) {
			err = suite.openAPIValidator.ValidateResponse(body, httpReq, httpResp.StatusCode, httpResp.Header)
			assert.NoError(suite.T(), err)
		})
	}
}

// TestFunctionalGetProjects runs the test suite
func TestFunctionalGetProjects(t *testing.T) {
	suite.Run(t, new(SuiteGetProjects))
}
