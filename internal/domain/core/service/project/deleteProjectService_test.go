package project

import (
	"fmt"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProjectService_Delete(t *testing.T) {

	tests := []struct {
		desc        string
		service     *DeleteProjectService
		id          string
		arrangeFunc func(*testing.T, *DeleteProjectService)
		assertFunc  func(*testing.T, *DeleteProjectService) bool
		err         error
	}{
		{
			desc:    "Testing an error deleting a project on the DeleteProjectService service when the project repository is not initialized",
			service: NewDeleteProjectService(nil, nil, logger.NewFakeLogger()),
			id:      "test-id",
			err:     fmt.Errorf(ErrProjectRepositoryNotInitialized),
		},
		{
			desc: "Testing an error deleting a project on the DeleteProjectService service when the project storage is not provided",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				nil,
				logger.NewFakeLogger(),
			),
			id:  "test-id",
			err: fmt.Errorf(ErrProjectStorageNotProvided),
		},
		{
			desc: "Testing an error deleting a project on the DeleteProjectService service when the project id is not provided",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			id: "",
			err: domainerror.NewProjectNotProvidedError(
				fmt.Errorf(ErrProjectIDNotProvided),
			),
		},
		{
			desc: "Testing an error deleting a project on the DeleteProjectService service when the project is not found",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			id: "test-id",
			arrangeFunc: func(t *testing.T, service *DeleteProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"test-id",
				).Return(
					nil,
					fmt.Errorf("project not found"),
				)
			},
			err: domainerror.NewProjectNotFoundError(
				fmt.Errorf("%s: %w", ErrFindingProject, fmt.Errorf("project not found")),
			),
		},
		{
			desc: "Testing an error deleting a project on the DeleteProjectService service when the project storage storer is not found",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			id: "test-id",
			arrangeFunc: func(t *testing.T, service *DeleteProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"test-id",
				).Return(
					&entity.Project{
						Name:    "test-id",
						Storage: "local",
					},
					nil,
				)

				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(
					nil,
				)
			},
			err: fmt.Errorf(ErrStorageHandlerNotFound),
		},
		{
			desc: "Testing an error deleting a project on the DeleteProjectService service when there is an error deleting the project from the repository",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			id: "test-id",
			arrangeFunc: func(t *testing.T, service *DeleteProjectService) {
				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"test-id",
				).Return(
					&entity.Project{
						Name:    "test-id",
						Storage: "local",
					},
					nil,
				)

				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(
					repository.NewMockProjectSourceCodeStorer(),
					nil,
				)

				service.repository.(*repository.MockProjectRepository).On(
					"Delete",
					"test-id",
				).Return(
					fmt.Errorf("error deleting project"),
				)
			},
			err: fmt.Errorf("%s: %w", ErrDeletingProject, fmt.Errorf("error deleting project")),
		},
		{
			desc: "Testing an error deleting a project on the DeleteProjectService service when there is an error deleting the project source code from the storage",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			id: "test-id",
			arrangeFunc: func(t *testing.T, service *DeleteProjectService) {
				projectSourceCodeStorer := repository.NewMockProjectSourceCodeStorer()

				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"test-id",
				).Return(
					&entity.Project{
						Name:    "test-id",
						Storage: "local",
					},
					nil,
				)

				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(
					projectSourceCodeStorer,
					nil,
				)

				service.repository.(*repository.MockProjectRepository).On(
					"Delete",
					"test-id",
				).Return(
					nil,
				)

				projectSourceCodeStorer.On(
					"Delete",
					&entity.Project{
						Name:    "test-id",
						Storage: "local",
					},
				).Return(
					fmt.Errorf("error deleting project source code"),
				)
			},
			err: fmt.Errorf("%s: %w", ErrDeletingProject, fmt.Errorf("error deleting project source code")),
		},
		{
			desc: "Testing successfully deleting a project on the DeleteProjectService service",
			service: NewDeleteProjectService(
				repository.NewMockProjectRepository(),
				repository.NewMockProjectSourceCodeStorageFactory(),
				logger.NewFakeLogger(),
			),
			id: "test-id",
			arrangeFunc: func(t *testing.T, service *DeleteProjectService) {
				projectSourceCodeStorer := repository.NewMockProjectSourceCodeStorer()

				service.repository.(*repository.MockProjectRepository).On(
					"Find",
					"test-id",
				).Return(
					&entity.Project{
						Name:    "test-id",
						Storage: "local",
					},
					nil,
				)

				service.storage.(*repository.MockProjectSourceCodeStorageFactory).On(
					"Get",
					"local",
				).Return(
					projectSourceCodeStorer,
					nil,
				)

				service.repository.(*repository.MockProjectRepository).On(
					"Delete",
					"test-id",
				).Return(
					nil,
				)

				projectSourceCodeStorer.On(
					"Delete",
					&entity.Project{
						Name:    "test-id",
						Storage: "local",
					},
				).Return(
					nil,
				)
			},
			assertFunc: func(t *testing.T, service *DeleteProjectService) bool {
				return service.repository.(*repository.MockProjectRepository).AssertExpectations(t) &&
					service.storage.(*repository.MockProjectSourceCodeStorageFactory).AssertExpectations(t)
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.service)
			}

			err := test.service.Delete(test.id)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Nil(t, err, "expected no error, got %v", err)
				assert.Nil(t, test.err, "no error received, but expected %v", test.err)

				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(t, test.service), "assertion function returned false")
				}
			}
		})
	}
}
