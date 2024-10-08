package mapper

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
)

// ProjectMapper is responsible for mapping project entity to response
type ProjectMapper struct{}

// NewProjectMapper creates a new project mapper
func NewProjectMapper() *ProjectMapper {
	return &ProjectMapper{}
}

// ToProjectResponse maps a project entity to a project response
func (m *ProjectMapper) ToProjectResponse(project *entity.Project) *response.ProjectResponse {

	if project == nil {
		return &response.ProjectResponse{}
	}

	return &response.ProjectResponse{
		Format:    project.Format,
		Name:      project.Name,
		Reference: project.Reference,
		Type:      project.Type,
	}
}
