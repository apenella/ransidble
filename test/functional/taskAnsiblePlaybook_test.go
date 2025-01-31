package functional

// import (
// 	"context"
// 	"fmt"
// 	nethttp "net/http"
// 	"testing"
// 	"time"

// 	"github.com/apenella/ransidble/internal/domain/core/entity"
// 	"github.com/apenella/ransidble/internal/domain/core/service/executor"
// 	ansibleexecutor "github.com/apenella/ransidble/internal/domain/core/service/executor"
// 	projectService "github.com/apenella/ransidble/internal/domain/core/service/project"
// 	taskService "github.com/apenella/ransidble/internal/domain/core/service/task"
// 	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
// 	serve "github.com/apenella/ransidble/internal/handler/cli/serve"
// 	"github.com/apenella/ransidble/internal/handler/http"
// 	projectHandler "github.com/apenella/ransidble/internal/handler/http/project"
// 	taskHandler "github.com/apenella/ransidble/internal/handler/http/task"
// 	"github.com/apenella/ransidble/internal/infrastructure/filesystem"
// 	"github.com/apenella/ransidble/internal/infrastructure/logger"
// 	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch"
// 	localprojectpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository"
// 	taskpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
// 	"github.com/apenella/ransidble/internal/infrastructure/tar"
// 	"github.com/apenella/ransidble/internal/infrastructure/unpack"
// 	"github.com/labstack/echo/v4"
// 	"github.com/spf13/afero"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// // TestSuite is the test suite for the HTTP server
// type TestSuite struct {
// 	listenAddress string
// 	server        *http.Server
// 	suite.Suite
// }

// // SetupSuite runs once before the suite starts running
// func (suite *TestSuite) SetupSuite() {

// 	log := logger.NewFakeLogger()
// 	afs := afero.NewOsFs()
// 	fs := filesystem.NewFilesystem(afs)

// 	// At this moment, the project repository loads the projects from the local storage. In the future, the plan is to have a database where you need to create a project before running it.
// 	projectsRepository := localprojectpersistence.NewLocalProjectRepository(
// 		afs,
// 		"../projects",
// 		log,
// 	)

// 	errLoadProjects := projectsRepository.LoadProjects()
// 	if errLoadProjects != nil {
// 		suite.T().Errorf("Error loading projects: %s", errLoadProjects)
// 		suite.T().FailNow()
// 		return
// 	}

// 	fetchFactory := fetch.NewFactory()
// 	fetchFactory.Register(
// 		entity.ProjectTypeLocal,
// 		fetch.NewLocalStorage(
// 			afs,
// 			log,
// 		),
// 	)

// 	unpackFactory := unpack.NewFactory()
// 	unpackFactory.Register(entity.ProjectFormatPlain, unpack.NewPlainFormat(
// 		afs,
// 		log,
// 	))

// 	tarExtractor := tar.NewTar(afs, log)
// 	unpackFactory.Register(entity.ProjectFormatTarGz, unpack.NewTarGzipFormat(
// 		afs,
// 		tarExtractor,
// 		log,
// 	))

// 	workspaceBuilder := workspace.NewBuilder(
// 		fs,
// 		fetchFactory,
// 		unpackFactory,
// 		projectsRepository,
// 		log,
// 	)

// 	ansiblePlaybookExecutor := ansibleexecutor.NewMockAnsiblePlaybookExecutor()
// 	ansiblePlaybookExecutor.On("Run", mock.Anything, mock.Anything).Return(nil)

// 	dispatcher := executor.NewDispatch(
// 		1,
// 		workspaceBuilder,
// 		ansiblePlaybookExecutor,
// 		log,
// 	)

// 	taskRepository := taskpersistence.NewMemoryTaskRepository(log)
// 	createTaskAnsiblePlaybookService := taskService.NewCreateTaskAnsiblePlaybookService(
// 		dispatcher,
// 		taskRepository,
// 		projectsRepository,
// 		log,
// 	)

// 	createTaskAnsiblePlaybookHandler := taskHandler.NewCreateTaskAnsiblePlaybookHandler(createTaskAnsiblePlaybookService, log)

// 	getTaskService := taskService.NewGetTaskService(taskRepository, log)
// 	getTaskHandler := taskHandler.NewGetTaskHandler(getTaskService, log)

// 	getProjectService := projectService.NewGetProjectService(projectsRepository, log)
// 	getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
// 	getProjectListHandler := projectHandler.NewGetProjectListHandler(getProjectService, log)

// 	router := echo.New()
// 	router.POST(serve.CreateTaskAnsiblePlaybookPath, createTaskAnsiblePlaybookHandler.Handle)
// 	router.GET(serve.GetTaskPath, getTaskHandler.Handle)
// 	router.GET(serve.GetProjectPath, getProjectHandler.Handle)
// 	router.GET(serve.GetProjectsPath, getProjectListHandler.Handle)

// 	suite.listenAddress = "0.0.0.0:8080"
// 	suite.server = http.NewServer(suite.listenAddress, router, log)
// }

// // TearDownSuite runs after all tests in this suite have run
// func (suite *TestSuite) TearDownSuite() {}

// // SetupTest runs before each test in the suite
// func (suite *TestSuite) SetupTest() {
// }

// // TearDownTest runs after each test in the suite
// func (suite *TestSuite) TearDownTest() {
// }

// // TestGetProjects is a functional test for the GetProjects endpoint
// func (suite *TestSuite) TestGetProjects() {
// 	var err error

// 	go func() {
// 		err = suite.server.Start(context.Background())
// 	}()

// 	defer suite.server.Stop()

// 	url := "http://" + suite.listenAddress + serve.GetProjectsPath
// 	fmt.Println(">>>>>", url)

// 	time.Sleep(3 * time.Second)

// 	resp, err := nethttp.Get(url)

// 	fmt.Println(">>>>>", resp)

// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), nethttp.StatusOK, resp.StatusCode)

// }

// // TestFunctionalTestSuite runs the test suite
// func TestFunctionalTestSuite(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }
