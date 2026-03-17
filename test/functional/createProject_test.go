package functional

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	nethttp "net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/request"
	projectService "github.com/apenella/ransidble/internal/domain/core/service/project"
	"github.com/apenella/ransidble/internal/handler/http"
	"github.com/apenella/ransidble/internal/handler/http/project"
	projectHandler "github.com/apenella/ransidble/internal/handler/http/project"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository/local"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/store"
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SuiteCreateProject is the test suite for the CreateProjectHandler
type SuiteCreateProject struct {
	listenAddress string
	router        *echo.Echo
	server        *http.Server

	suite.Suite
}

// SetupSuite runs once before the suite starts running
func (suite *SuiteCreateProject) SetupSuite() {
	suite.listenAddress = "0.0.0.0:8080"
}

// SetupTest runs before each test in the suite
func (suite *SuiteCreateProject) SetupTest() {
	suite.router = echo.New()
	suite.server = http.NewServer(suite.listenAddress, suite.router, logger.NewLogger())
}

// TearDownSuite runs after all tests in this suite have run
// func (suite *SuiteCreateProject) TearDownSuite() {
// 	suite.server.Stop()
// }

// TearDownTest runs after each test in the suite
func (suite *SuiteCreateProject) TearDownTest() {
	suite.server.Stop()
}

// TestCreateProject is a functional test for the CreateProject endpoint
func (suite *SuiteCreateProject) TestCreateProject() {
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
		desc                   string
		arrangeTest            func(*SuiteCreateProject)
		prepareInputParameters func(description string) (*InputFunctionalTest, error)
	}{
		{
			desc: "Testing functional behaviour to create a project and request is completed successfully and returns a StatusCreated",
			prepareInputParameters: func(description string) (*InputFunctionalTest, error) {
				var bodyBuffer bytes.Buffer
				var expectedBody string

				method := "POST"
				url := "http://localhost:8080/projects/project-1"
				expectedStatusCode := nethttp.StatusCreated
				expectedHttpHeaders := map[string]string{
					"Location": "/projects/project-1",
				}

				input := &InputFunctionalTest{
					desc:                description,
					method:              method,
					url:                 url,
					expectedStatusCode:  expectedStatusCode,
					expectedHttpHeaders: expectedHttpHeaders,
					expectedBody:        expectedBody,
				}

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: entity.ProjectTypeLocal,
				}
				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					return nil, err
				}

				multipartWriter := multipart.NewWriter(&bodyBuffer)
				defer multipartWriter.Close()

				multipartWriter.WriteField(project.RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				part, err := multipartWriter.CreateFormFile(project.RequestFormProjectFileFieldeName, "project-1.tar.gz")
				if err != nil {
					return nil, err
				}
				projectContentFile := strings.NewReader("project-content")
				_, err = io.Copy(part, projectContentFile)
				if err != nil {
					return nil, err
				}

				headers := map[string]string{
					echo.HeaderContentType: fmt.Sprintf("multipart/form-data; boundary=%s", multipartWriter.Boundary()),
				}

				input.headers = headers
				input.parameters = io.NopCloser(&bodyBuffer)

				return input, nil
			},
			arrangeTest: func(suite *SuiteCreateProject) {
				// arrangeTest function prepares the handler that receives the request during the test.

				log := logger.NewLogger()

				roFsBase := afero.NewReadOnlyFs(afero.NewOsFs())
				rwFs := afero.NewCopyOnWriteFs(roFsBase, afero.NewMemMapFs())

				projectsRepository := local.NewDatabaseDriver(
					rwFs,
					filepath.Join("..", "fixtures", "functional-create-project"),
					log,
				)

				storeFactory := store.NewFactory()
				localStorageStore := store.NewLocalStorage(
					rwFs,
					filepath.Join("..", "fixtures", "functional-create-project"),
					log,
				)
				storeFactory.Register(entity.ProjectTypeLocal, localStorageStore)

				createProjectService := projectService.NewCreateProjectService(
					projectsRepository,
					storeFactory,
					log,
				)

				localStorageStore.Initialize()
				projectsRepository.Initialize()

				createProjectHander := projectHandler.NewCreateProjectHandler(createProjectService, log)

				suite.router.POST(http.CreateProjectPath, createProjectHander.Handle)
			},
		},
		{
			desc: "Testing functional behaviour to create a project and request metadata is not properly provided and is returned a StatusBadRequest",
			prepareInputParameters: func(description string) (*InputFunctionalTest, error) {
				var bodyBuffer bytes.Buffer

				method := "POST"
				url := "http://localhost:8080/projects/project-1"
				expectedStatusCode := nethttp.StatusBadRequest
				expectedHttpHeaders := map[string]string{}
				expectedBody := "{\"id\":\"\",\"error\":\"project metadata not provided in the request\",\"status\":400}"
				input := &InputFunctionalTest{
					desc:                description,
					method:              method,
					url:                 url,
					expectedStatusCode:  expectedStatusCode,
					expectedHttpHeaders: expectedHttpHeaders,
					expectedBody:        expectedBody,
				}

				multipartWriter := multipart.NewWriter(&bodyBuffer)
				defer multipartWriter.Close()

				part, err := multipartWriter.CreateFormFile(project.RequestFormProjectFileFieldeName, "project-1.tar.gz")
				if err != nil {
					return nil, err
				}
				projectContentFile := strings.NewReader("project-content")
				_, err = io.Copy(part, projectContentFile)
				if err != nil {
					return nil, err
				}

				headers := map[string]string{
					echo.HeaderContentType: fmt.Sprintf("multipart/form-data; boundary=%s", multipartWriter.Boundary()),
				}

				input.headers = headers
				input.parameters = io.NopCloser(&bodyBuffer)

				return input, nil
			},
			arrangeTest: func(suite *SuiteCreateProject) {
				// arrangeTest function prepares the handler that receives the request during the test.

				log := logger.NewLogger()

				roFsBase := afero.NewReadOnlyFs(afero.NewOsFs())
				rwFs := afero.NewCopyOnWriteFs(roFsBase, afero.NewMemMapFs())

				projectsRepository := local.NewDatabaseDriver(
					rwFs,
					filepath.Join("..", "fixtures", "functional-create-project"),
					log,
				)

				storeFactory := store.NewFactory()
				localStorageStore := store.NewLocalStorage(
					rwFs,
					filepath.Join("..", "fixtures", "functional-create-project"),
					log,
				)
				storeFactory.Register(entity.ProjectTypeLocal, localStorageStore)

				createProjectService := projectService.NewCreateProjectService(
					projectsRepository,
					storeFactory,
					log,
				)

				localStorageStore.Initialize()
				projectsRepository.Initialize()

				createProjectHander := projectHandler.NewCreateProjectHandler(createProjectService, log)

				suite.router.POST(http.CreateProjectPath, createProjectHander.Handle)
			},
		},

		{
			desc: "Testing functional behaviour to create a project and request file is not properly provided and is returned a StatusBadRequest",
			prepareInputParameters: func(description string) (*InputFunctionalTest, error) {
				var bodyBuffer bytes.Buffer

				method := "POST"
				url := "http://localhost:8080/projects/project-1"
				expectedStatusCode := nethttp.StatusBadRequest
				expectedHttpHeaders := map[string]string{}
				expectedBody := "{\"id\":\"\",\"error\":\"error reading project file field: http: no such file\",\"status\":400}"
				input := &InputFunctionalTest{
					desc:                description,
					method:              method,
					url:                 url,
					expectedStatusCode:  expectedStatusCode,
					expectedHttpHeaders: expectedHttpHeaders,
					expectedBody:        expectedBody,
				}

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: entity.ProjectTypeLocal,
				}
				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					return nil, err
				}

				multipartWriter := multipart.NewWriter(&bodyBuffer)
				defer multipartWriter.Close()

				multipartWriter.WriteField(project.RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				headers := map[string]string{
					echo.HeaderContentType: fmt.Sprintf("multipart/form-data; boundary=%s", multipartWriter.Boundary()),
				}

				input.headers = headers
				input.parameters = io.NopCloser(&bodyBuffer)

				return input, nil
			},
			arrangeTest: func(suite *SuiteCreateProject) {
				// arrangeTest function prepares the handler that receives the request during the test.

				log := logger.NewLogger()

				roFsBase := afero.NewReadOnlyFs(afero.NewOsFs())
				rwFs := afero.NewCopyOnWriteFs(roFsBase, afero.NewMemMapFs())

				projectsRepository := local.NewDatabaseDriver(
					rwFs,
					filepath.Join("..", "fixtures", "functional-create-project"),
					log,
				)

				storeFactory := store.NewFactory()
				localStorageStore := store.NewLocalStorage(
					rwFs,
					filepath.Join("..", "fixtures", "functional-create-project"),
					log,
				)
				storeFactory.Register(entity.ProjectTypeLocal, localStorageStore)

				createProjectService := projectService.NewCreateProjectService(
					projectsRepository,
					storeFactory,
					log,
				)

				localStorageStore.Initialize()
				projectsRepository.Initialize()

				createProjectHander := projectHandler.NewCreateProjectHandler(createProjectService, log)

				suite.router.POST(http.CreateProjectPath, createProjectHander.Handle)
			},
		},
	}

	for _, test := range tests {

		suite.T().Log(test.desc)

		if test.arrangeTest != nil {
			test.arrangeTest(suite)
		}

		input, err := test.prepareInputParameters(test.desc)
		assert.NoError(suite.T(), err)

		err = actAndAssert(suite.T(), input)
		assert.NoError(suite.T(), err)
	}

}

// TestCreateProjectVersion is a functional test for the CreateProjectVersion endpoint
func TestCreateProject(t *testing.T) {
	suite.Run(t, new(SuiteCreateProject))
}
