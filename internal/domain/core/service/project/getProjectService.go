package project

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

var (
	// ErrFindingProject represents an error when a proejct is not found
	ErrFindingProject = fmt.Errorf("error finding project")
	// ErrRepositoryNotInitialized represents an error when the store is not initialized
	ErrRepositoryNotInitialized = fmt.Errorf("project repository not initialized")
	// ErrProjectIDNotProvided represents an error when the task id is not provided
	ErrProjectIDNotProvided = fmt.Errorf("project id not provided")
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
		p.logger.Error(ErrRepositoryNotInitialized.Error(), map[string]interface{}{
			"component":  "GetProjectService.GetProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": id,
		})
		return nil, ErrRepositoryNotInitialized
	}

	if id == "" {
		p.logger.Error(ErrProjectIDNotProvided.Error(), map[string]interface{}{
			"component": "GetProjectService.GetProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return nil, domainerror.NewProjectNotProvidedError(ErrProjectIDNotProvided)
	}

	project, err := p.repository.Find(id)
	if err != nil {
		p.logger.Error("%s: %s", ErrFindingProject.Error(), err.Error(), map[string]interface{}{
			"component":  "GetProjectService.GetProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": id,
		})
		return nil, domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s: %s", ErrFindingProject.Error(), err.Error()),
		)
	}

	return project, nil
}

// GetProjectsList returns all projects
func (p *GetProjectService) GetProjectsList() ([]*entity.Project, error) {

	if p.repository == nil {
		p.logger.Error(ErrRepositoryNotInitialized.Error(), map[string]interface{}{
			"component": "GetProjectService.GetProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return nil, ErrRepositoryNotInitialized
	}

	projects, err := p.repository.FindAll()
	if err != nil {
		p.logger.Error("%s: %s", ErrFindingProject.Error(), err.Error(), map[string]interface{}{
			"component": "GetProjectService.GetProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return nil, domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s: %s", ErrFindingProject.Error(), err.Error()),
		)
	}

	return projects, nil
}
