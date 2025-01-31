package functional

import (
	"context"
	"fmt"
	"io"
	"net"
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

// TestSuite is the test suite for the HTTP server
type TestSuite struct {
	getProjectListHandler *projectHandler.GetProjectListHandler
	getProjectHandler     *projectHandler.GetProjectHandler
	listenAddress         string

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *TestSuite) SetupSuite() {

	suite.listenAddress = "0.0.0.0:8080"

	log := logger.NewFakeLogger()
	afs := afero.NewOsFs()

	// At this moment, the project repository loads the projects from the local storage. In the future, the plan is to have a database where you need to create a project before running it.
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
	suite.getProjectHandler = projectHandler.NewGetProjectHandler(getProjectService, log)
	suite.getProjectListHandler = projectHandler.NewGetProjectListHandler(getProjectService, log)
}

// TearDownSuite runs after all tests in this suite have run
func (suite *TestSuite) TearDownSuite() {}

// SetupTest runs before each test in the suite
func (suite *TestSuite) SetupTest() {
}

// TearDownTest runs after each test in the suite
func (suite *TestSuite) TearDownTest() {
}

// TestGetProjects is a functional test for the GetProjects endpoint
func (suite *TestSuite) TestGetProjects() {
	var err error

	router := echo.New()
	router.GET(serve.GetProjectsPath, suite.getProjectListHandler.Handle)

	server := http.NewServer(suite.listenAddress, router, logger.NewFakeLogger())

	go func() {
		err = server.Start(context.Background())
	}()
	defer server.Stop()

	url := "http://" + suite.listenAddress + serve.GetProjectsPath

	for i := 0; i < 5; i++ {
		conn, err := net.DialTimeout("tcp", suite.listenAddress, 1*time.Second)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}
	resp, err := nethttp.Get(url)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Errorf("Error reading response body: %s", err)
		suite.T().FailNow()
		return
	}
	defer resp.Body.Close()

	fmt.Println(">>>>>", string(body))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), nethttp.StatusOK, resp.StatusCode)

}

// TestGetProjects is a functional test for the GetProjects endpoint
func (suite *TestSuite) TestGetProjectProject1() {
	var err error

	router := echo.New()
	router.GET(serve.GetProjectPath, suite.getProjectHandler.Handle)

	server := http.NewServer(suite.listenAddress, router, logger.NewFakeLogger())

	go func() {
		err = server.Start(context.Background())
	}()
	defer server.Stop()

	url := "http://" + suite.listenAddress + serve.GetProjectsPath + "/project-1"

	for i := 0; i < 5; i++ {
		conn, err := net.DialTimeout("tcp", suite.listenAddress, 1*time.Second)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}

	resp, err := nethttp.Get(url)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Errorf("Error reading response body: %s", err)
		suite.T().FailNow()
		return
	}
	defer resp.Body.Close()

	fmt.Println(">>>>>", string(body))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), nethttp.StatusOK, resp.StatusCode)

}

// TestFunctionalTestSuite runs the test suite
func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
