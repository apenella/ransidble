package mapper

import (
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/stretchr/testify/assert"
)

// TestToProjectResponse maps a project entity to a project response
func TestToProjectResponse(t *testing.T) {
	tests := []struct {
		desc     string
		project  *entity.Project
		mapper   *ProjectMapper
		expected *response.ProjectResponse
	}{
		{
			desc: "Testing project mapping",
			project: &entity.Project{
				Format:    "project-format",
				Name:      "project-name",
				Reference: "project-reference",
				Storage:   "project-storage",
			},
			expected: &response.ProjectResponse{
				Format:    "project-format",
				Name:      "project-name",
				Reference: "project-reference",
				Storage:   "project-storage",
			},
			mapper: NewProjectMapper(),
		},
		{
			desc:     "Testing project mapping with empty project",
			project:  &entity.Project{},
			expected: &response.ProjectResponse{},
			mapper:   NewProjectMapper(),
		},
		{
			desc:     "Testing project mapping with nil project",
			project:  nil,
			expected: &response.ProjectResponse{},
			mapper:   NewProjectMapper(),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			res := test.mapper.ToProjectResponse(test.project)
			assert.Equal(t, test.expected, res)

		})
	}
}
