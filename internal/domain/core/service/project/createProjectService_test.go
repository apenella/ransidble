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
		arrangeFunc          func(*testing.T, *CreateProjectService)
		assertFunc           func(*testing.T, *CreateProjectService) bool
		desc                 string
		err                  error
		format               string
		projectContentReader io.Reader
		projectID            string
		projectVersion       string
		service              *CreateProjectService
		storage              string
	}{
		{
			desc:                 "Testing create a project on the CreateProjectService providing a specific version",
			format:               "targz",
			storage:              "local",
			projectID:            "project-id",
			projectVersion:       "v1.0.0",
			projectContentReader: fileReader,
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
					"project-id",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project-id",
					&entity.Project{
						Name:      "project-id",
						Version:   "v1.0.0",
						Format:    "targz",
						Storage:   "local",
						Reference: "project-id.tar.gz",
					},
				).Return(nil)

				projectSourceCodeStorer.On(
					"Store",
					&entity.Project{
						Name:      "project-id",
						Version:   "v1.0.0",
						Format:    "targz",
						Storage:   "local",
						Reference: "project-id.tar.gz",
					},
					fileReader,
				).Return(nil)
			},
			assertFunc: func(t *testing.T, service *CreateProjectService) bool {
				return service.repository.(*repository.MockProjectRepository).AssertExpectations(t)
			},
		},
		{
			desc:                 "Testing create a project on the CreateProjectService without providing a version",
			format:               "targz",
			storage:              "local",
			projectID:            "project-id",
			projectVersion:       "",
			projectContentReader: fileReader,
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
					"project-id",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project-id",
					&entity.Project{
						Name:      "project-id",
						Version:   "latest",
						Format:    "targz",
						Storage:   "local",
						Reference: "project-id.tar.gz",
					},
				).Return(nil)

				projectSourceCodeStorer.On(
					"Store",
					&entity.Project{
						Name:      "project-id",
						Version:   "latest",
						Format:    "targz",
						Storage:   "local",
						Reference: "project-id.tar.gz",
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
			projectID:            "project-id",
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
			projectID:            "project-id",
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
			projectID:            "project",
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
			projectID:            "",
			err: domainerror.NewProjectIDNotProvidedError(
				fmt.Errorf(ErrProjectIDNotProvided),
			),
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
			projectID:            "project-id",
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
			projectID:            "project-id",
			err:                  fmt.Errorf(ErrProjectRepositoryNotInitialized),
			service: NewCreateProjectService(
				nil,
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
			projectID:            "project-id",
			err: fmt.Errorf(
				"%s: %s",
				ErrProjectFormatNotSupported,
				"invalid format: non-supported-format",
			),
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
			projectID:            "project-id",
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
					"project-id",
				).Return(&entity.Project{
					Name:      "project-id",
					Format:    "plain",
					Storage:   "local",
					Reference: "project-id.tar.gz",
				}, nil)
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storage in not supported",
			format:               "plain",
			storage:              "non-supported-storage",
			projectContentReader: fileReader,
			projectID:            "project-id",
			err: fmt.Errorf(
				"%s: %s",
				ErrProjectStorageNotSupported, "invalid storage type: non-supported-storage",
			),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storage handler is not found",
			format:               "plain",
			storage:              "local",
			projectContentReader: fileReader,
			projectID:            "project-id",
			err:                  fmt.Errorf(ErrStorageHandlerNotFound),
			service: NewCreateProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *CreateProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"project-id",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(nil)
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storing a project to the repository fails",
			format:               "targz",
			storage:              "local",
			projectContentReader: fileReader,
			projectID:            "project-id",
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
					"project-id",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project-id",
					&entity.Project{
						Format:    "targz",
						Name:      "project-id",
						Reference: "project-id.tar.gz",
						Storage:   "local",
						Version:   "latest",
					},
				).Return(fmt.Errorf("storing project fails"))
			},
		},
		{
			desc:                 "Testing an error creating a project on the CreateProjectService service when storing a project to the persistent storage fails",
			format:               "targz",
			storage:              "local",
			projectContentReader: fileReader,
			projectID:            "project-id",
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
					"project-id",
				).Return(nil, nil)
				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(projectSourceCodeStorer)
				service.repository.(*repository.MockProjectRepository).On(
					"SafeStore",
					"project-id",
					&entity.Project{
						Format:    "targz",
						Name:      "project-id",
						Reference: "project-id.tar.gz",
						Storage:   "local",
						Version:   "latest",
					},
				).Return(nil)

				projectSourceCodeStorer.On(
					"Store",
					&entity.Project{
						Format:    "targz",
						Name:      "project-id",
						Reference: "project-id.tar.gz",
						Storage:   "local",
						Version:   "latest",
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

			err := test.service.Create(test.format, test.storage, test.projectID, test.projectVersion, test.projectContentReader)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Nil(t, err, "unexpected error received")
				assert.Nil(t, test.err, "no error received when an error was expected")

				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(t, test.service))
				}
			}
		})
	}
}
