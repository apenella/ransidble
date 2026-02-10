package project

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestCreateProjectService_Create(t *testing.T) {

	fileReader := io.NopCloser(strings.NewReader("content for testing"))

	tests := []struct {
		desc                 string
		format               string
		storage              string
		fileName             string
		projectContentReader io.Reader
		expectedProjectID    string
		err                  error
		service              *CreateProjectService
		arrangeFunc          func(*testing.T, *CreateProjectService)
		assertFunc           func(*testing.T, *CreateProjectService) bool
	}{
		{
			desc:                 "Testing create a project on the CreateProjectService",
			format:               "plain",
			storage:              "local",
			fileName:             "project.tar.gz",
			projectContentReader: fileReader,
			expectedProjectID:    "project",
			err:                  nil,
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				projectSourceCodeStorer := repository.NewMockProjectSourceCodeStorer()

				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project",
					&entity.Project{
						Name:      "project",
						Format:    "plain",
						Storage:   "local",
						Reference: "project.tar.gz",
					},
				).Return(nil)

				projectSourceCodeStorer.On(
					"Store",
					&entity.Project{
						Name:      "project",
						Format:    "plain",
						Storage:   "local",
						Reference: "project.tar.gz",
					},
					fileReader,
				).Return(nil)
			},
			assertFunc: func(t *testing.T, service *CreateProjectService) bool {
				return service.repository.(*repository.MockProjectRepository).AssertExpectations(t)
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when the format is not provided",
			format:               "",
			storage:              "local",
			fileName:             "project.tar.gz",
			projectContentReader: fileReader,
			err:                  fmt.Errorf(ErrProjectFormatNotProvided),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when the storage is not provided",
			format:               "plain",
			storage:              "",
			fileName:             "project.tar.gz",
			projectContentReader: fileReader,
			err:                  fmt.Errorf(ErrProjectStorageNotProvided),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when the file is not provided",
			format:               "plain",
			storage:              "local",
			projectContentReader: nil,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf(ErrProjectContentReaderNotProvided),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when the file name is not provided",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "",
			err:                  fmt.Errorf(ErrFileNameNotProvided),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when the storage handler is not initialized",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf(ErrStorageHandlerNotInitialized),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				nil,
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when project repository is not initialized",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf(ErrProjectRepositoryNotInitialized),
			service: NewCreateProjectService(
				nil,
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when project name has an non supported extension",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.non-supported-extension",
			err:                  fmt.Errorf("%s: %s", ErrProjectFileExtensionNotSupported, "file project.non-supported-extension extension not supported"),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when format is not supported",
			format:               "non-supported-format",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf("%s: %s", ErrProjectFormatNotSupported, "invalid format: non-supported-format"),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when project already exists",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err: domainerror.NewProjectAlreadyExistsError(
				fmt.Errorf(ErrProjectAlreadyExists),
			),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project",
				).Return(&entity.Project{
					Name:      "project",
					Format:    "plain",
					Storage:   "local",
					Reference: "project.tar.gz",
				}, nil)
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storage in not supported",
			format:               "plain",
			storage:              "non-supported-storage",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf("%s: %s", ErrProjectStorageNotSupported, "invalid storage type: non-supported-storage"),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project",
				).Return(nil, nil)
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storage handler is not found",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf(ErrStorageHandlerNotFound),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(nil)
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storing a project to the repository fails",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf("%s: %s", ErrStoringProject, "storing project fails"),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				projectSourceCodeStorer := repository.NewMockProjectSourceCodeStorer()

				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project",
					&entity.Project{
						Name:      "project",
						Format:    "plain",
						Storage:   "local",
						Reference: "project.tar.gz",
					},
				).Return(fmt.Errorf("storing project fails"))
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storing a project to the persistent storage fails",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			fileName:             "project.tar.gz",
			err:                  fmt.Errorf("%s: %s", ErrStoringProject, "storing project fails"),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				projectSourceCodeStorer := repository.NewMockProjectSourceCodeStorer()

				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project",
					&entity.Project{
						Name:      "project",
						Format:    "plain",
						Storage:   "local",
						Reference: "project.tar.gz",
					},
				).Return(nil)

				projectSourceCodeStorer.On(
					"Store",
					&entity.Project{
						Name:      "project",
						Format:    "plain",
						Storage:   "local",
						Reference: "project.tar.gz",
					},
					fileReader,
				).Return(fmt.Errorf("storing project fails"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.service)
			}

			id, err := test.service.Create(test.format, test.storage, test.fileName, test.projectContentReader)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Nil(t, err, "unexpected error received")
				assert.Nil(t, test.err, "no error received when an error was expected")

				if test.assertFunc != nil {
					assert.Equal(t, test.expectedProjectID, id, "unexpected project ID returned")
					assert.True(t, test.assertFunc(t, test.service))
				}
			}
		})
	}
}

func TestExtractProjectName(t *testing.T) {
	tests := []struct {
		desc     string
		fileName string
		expected string
	}{
		{
			desc:     "Testing extract project name from a file name",
			fileName: "project.tar.gz",
			expected: "project",
		},
		{
			desc:     "Testing extract project name from a file name with a single extension",
			fileName: "project.tar",
			expected: "project",
		},
		{
			desc:     "Testing extract project name from a file name with no extension",
			fileName: "project",
			expected: "project",
		},
		{
			desc:     "Testing extract project name from a blank file name",
			fileName: "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			projectName := extractProjectName(test.fileName)
			assert.Equal(t, test.expected, projectName)
		})
	}
}
