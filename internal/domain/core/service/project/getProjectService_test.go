package project

import (
	"errors"
	"fmt"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetProject(t *testing.T) {

	tests := []struct {
		desc        string
		id          string
		err         error
		expected    *entity.Project
		service     *GetProjectService
		arrangeFunc func(*testing.T, *GetProjectService)
	}{
		{
			desc: "Testing getting a project on the GetProjectService",
			id:   "project-id",
			err:  errors.New(""),
			expected: &entity.Project{
				Name:      "project-id",
				Reference: "project-id",
				Format:    "plain",
				Storage:   "local",
			},
			service: NewGetProjectService(
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *GetProjectService) {
				service.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(&entity.Project{
					Name:      "project-id",
					Reference: "project-id",
					Format:    "plain",
					Storage:   "local",
				}, nil)
			},
		},
		{
			desc:     "Testing error getting a project on the GetProjectService having a nil project repository",
			id:       "project-id",
			err:      fmt.Errorf(ErrProjectRepositoryNotInitialized),
			expected: nil,
			service: NewGetProjectService(
				nil,
				logger.NewFakeLogger(),
			),
			arrangeFunc: nil,
		},
		{
			desc: "Testing error getting a project on the GetProjectService having an empty project id",
			id:   "",
			err: domainerror.NewProjectNotProvidedError(
				fmt.Errorf(ErrProjectIDNotProvided),
			),
			expected: nil,
			service: &GetProjectService{
				repository: repository.NewMockProjectRepository(),
				logger:     logger.NewFakeLogger(),
			},
			arrangeFunc: nil,
		},
		{
			desc: "Testing error getting a project on the GetProjectService having an error on find project into the repository",
			id:   "project-id",
			err: domainerror.NewProjectNotFoundError(
				fmt.Errorf("%s: %w", ErrFindingProject, errors.New("error finding project")),
			),
			expected: nil,
			service: &GetProjectService{
				repository: repository.NewMockProjectRepository(),
				logger:     logger.NewFakeLogger(),
			},
			arrangeFunc: func(t *testing.T, service *GetProjectService) {
				service.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(nil, errors.New("error finding project"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.service)
			}

			project, err := test.service.GetProject(test.id)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, project)
			}
		})
	}
}

func TestGetProjectsList(t *testing.T) {

	tests := []struct {
		desc        string
		err         error
		expected    []*entity.Project
		service     *GetProjectService
		arrangeFunc func(*testing.T, *GetProjectService)
	}{
		{
			desc: "Testing getting a project list on the GetProjectService",
			err:  errors.New(""),
			expected: []*entity.Project{
				{
					Name:      "project-id",
					Reference: "project-id",
					Format:    "plain",
					Storage:   "local",
				},
			},
			service: NewGetProjectService(
				repository.NewMockProjectRepository(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, service *GetProjectService) {
				service.repository.(*repository.MockProjectRepository).On("FindAll").Return([]*entity.Project{
					{
						Name:      "project-id",
						Reference: "project-id",
						Format:    "plain",
						Storage:   "local",
					},
				}, nil)
			},
		},
		{
			desc:     "Testing error getting a project list on the GetProjectService having a nil project repository",
			err:      fmt.Errorf(ErrProjectRepositoryNotInitialized),
			expected: nil,
			service: NewGetProjectService(
				nil,
				logger.NewFakeLogger(),
			),
			arrangeFunc: nil,
		},
		{
			desc: "Testing error getting a project list on the GetProjectService having an error on find project into the repository",
			err: domainerror.NewProjectNotFoundError(
				fmt.Errorf("%s: %s", ErrFindingProject, errors.New("error finding project")),
			),
			expected: nil,
			service: &GetProjectService{
				repository: repository.NewMockProjectRepository(),
				logger:     logger.NewFakeLogger(),
			},
			arrangeFunc: func(t *testing.T, service *GetProjectService) {
				service.repository.(*repository.MockProjectRepository).On("FindAll").Return(nil, errors.New("error finding project"))
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

			projects, err := test.service.GetProjectsList()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, projects)
			}
		})
	}

}
