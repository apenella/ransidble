package functional

import (
	"context"
	"fmt"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	taskService "github.com/apenella/ransidble/internal/domain/core/service/task"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/handler/cli/serve"
	"github.com/apenella/ransidble/internal/handler/http"
	taskHandler "github.com/apenella/ransidble/internal/handler/http/task"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SuiteGetTask is the test suite for the HTTP server
type SuiteGetTask struct {
	listenAddress string
	router        *echo.Echo
	server        *http.Server

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *SuiteGetTask) SetupSuite() {
	suite.listenAddress = "0.0.0.0:8080"
}

// SetupTest runs before each test
func (suite *SuiteGetTask) SetupTest() {
	suite.router = echo.New()
	suite.server = http.NewServer(suite.listenAddress, suite.router, logger.NewFakeLogger())
}

// TearDownTest runs after the suite ends
func (suite *SuiteGetTask) TearDownTest() {
	suite.server.Stop()
}

// TestGetTask tests the GetTask method
func (suite *SuiteGetTask) TestGetTask() {
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
		expectedBody       string
		arrangeTest        func()
	}{
		{
			desc:               "Testing a request to get an existing task and return a StatusOK",
			method:             "GET",
			url:                "http://" + suite.listenAddress + "/tasks/project-1",
			expectedStatusCode: nethttp.StatusOK,
			expectedBody:       "{\"command\":\"ansible-playbook\",\"completed_at\":\"0000-01-01T01:01:01\",\"created_at\":\"0000-01-01T01:01:01\",\"executed_at\":\"0000-01-01T01:01:01\",\"id\":\"1\",\"parameters\":{\"playbooks\":[\"playbook.yml\"],\"inventory\":\"inventory.yml\"},\"project_id\":\"project-1\",\"status\":\"ACCEPTED\"}",
			arrangeTest: func() {
				// the task repository is mocked and returns a valid task
				repository := repository.NewMockTaskRepository()
				repository.On("Find", "project-1").Return(&entity.Task{
					ID:        "1",
					ProjectID: "project-1",
					Command:   entity.AnsiblePlaybookCommand,
					Parameters: &entity.AnsiblePlaybookParameters{
						Playbooks: []string{"playbook.yml"},
						Inventory: "inventory.yml",
					},
					CompletedAt: "0000-01-01T01:01:01",
					CreatedAt:   "0000-01-01T01:01:01",
					ExecutedAt:  "0000-01-01T01:01:01",
					Status:      entity.ACCEPTED,
				}, nil)

				arrangeGetTaskRouter(suite.router, serve.GetTaskPath, repository)
			},
		},
		{
			desc:               "Testing a request to get a non-existing task and return a StatusNotFound",
			method:             "GET",
			url:                "http://" + suite.listenAddress + "/tasks/project-1",
			expectedStatusCode: nethttp.StatusNotFound,
			expectedBody:       "{\"id\":\"\",\"error\":\"error getting task: error finding task project-1: task not found\"}",
			arrangeTest: func() {
				// the task repository is mocked and returns a task not found error
				repository := repository.NewMockTaskRepository()
				repository.On("Find", "project-1").Return(&entity.Task{}, domainerror.NewTaskNotFoundError(fmt.Errorf("task not found")))

				arrangeGetTaskRouter(suite.router, serve.GetTaskPath, repository)
			},
		},
	}

	for _, test := range tests {

		if test.arrangeTest != nil {
			test.arrangeTest()
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

// TestSuiteGetTask runs the test suite
func TestSuiteGetTask(t *testing.T) {
	suite.Run(t, new(SuiteGetTask))
}

func arrangeGetTaskRouter(router *echo.Echo, path string, repository *repository.MockTaskRepository) {

	getTaskService := taskService.NewGetTaskService(repository, logger.NewFakeLogger())
	getTaskHandler := taskHandler.NewGetTaskHandler(getTaskService, logger.NewFakeLogger())

	router.GET(path, getTaskHandler.Handle)
}
