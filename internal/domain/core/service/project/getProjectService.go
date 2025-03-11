package project

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

// GetProjectService is a service to get a project
type GetProjectService struct {
	repository repository.ProjectRepository
	logger     repository.Logger
}

// NewGetProjectService creates a new GetProjectService
func NewGetProjectService(repository repository.ProjectRepository, logger repository.Logger) *GetProjectService {
	return &GetProjectService{
		repository: repository,
		logger:     logger,
	}
}

// GetProject returns a project by its id
func (p *GetProjectService) GetProject(id string) (*entity.Project, error) {

	if p.repository == nil {
		p.logger.Error(ErrProjectRepositoryNotInitialized, map[string]interface{}{
			"component":  "GetProjectService.GetProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": id,
		})
		return nil, fmt.Errorf(ErrProjectRepositoryNotInitialized)
	}

	if id == "" {
		p.logger.Error(ErrProjectIDNotProvided, map[string]interface{}{
			"component": "GetProjectService.GetProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return nil, domainerror.NewProjectNotProvidedError(
			fmt.Errorf(ErrProjectIDNotProvided),
		)
	}

	project, err := p.repository.Find(id)
	if err != nil {
		p.logger.Error("%s: %s", ErrFindingProject, err.Error(), map[string]interface{}{
			"component":  "GetProjectService.GetProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": id,
		})
		return nil, domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s: %w", ErrFindingProject, err),
		)
	}

	return project, nil
}

// GetProjectsList returns all projects
func (p *GetProjectService) GetProjectsList() ([]*entity.Project, error) {

	if p.repository == nil {
		p.logger.Error(ErrProjectRepositoryNotInitialized, map[string]interface{}{
			"component": "GetProjectService.GetProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return nil, fmt.Errorf(ErrProjectRepositoryNotInitialized)
	}

	projects, err := p.repository.FindAll()
	if err != nil {
		p.logger.Error("%s: %s", ErrFindingProject, err.Error(), map[string]interface{}{
			"component": "GetProjectService.GetProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return nil, domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s: %s", ErrFindingProject, err.Error()),
		)
	}

	return projects, nil
}
