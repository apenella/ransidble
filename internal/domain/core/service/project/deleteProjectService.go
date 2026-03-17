package project

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
)

// DeleteProjectService is a service that handles the deletion of projects
type DeleteProjectService struct {
	repository repository.ProjectRepository
	storage    repository.SourceCodeStorageFactory
	logger     repository.Logger
}

// Ensure DeleteProjectService implements the DeleteProjectServicer interface
var _ service.DeleteProjectServicer = (*DeleteProjectService)(nil)

// NewDeleteProjectService creates a new instance of DeleteProjectService
func NewDeleteProjectService(repository repository.ProjectRepository, storage repository.SourceCodeStorageFactory, logger repository.Logger) *DeleteProjectService {
	return &DeleteProjectService{
		repository: repository,
		storage:    storage,
		logger:     logger,
	}
}

// Delete deletes a project by its id
func (s *DeleteProjectService) Delete(projectID string) error {

	var project *entity.Project
	var err error
	var storer repository.SourceCodeStorer

	if s.repository == nil {
		s.logger.Error(ErrProjectRepositoryNotInitialized, map[string]interface{}{
			"component":  "DeleteProjectService.DeleteProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": projectID,
		})
		return fmt.Errorf(ErrProjectRepositoryNotInitialized)
	}

	if s.storage == nil {
		s.logger.Error(ErrProjectStorageNotProvided, map[string]interface{}{
			"component":  "DeleteProjectService.DeleteProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": projectID,
		})
		return fmt.Errorf(ErrProjectStorageNotProvided)
	}

	if projectID == "" {
		s.logger.Error(ErrProjectIDNotProvided, map[string]interface{}{
			"component": "DeleteProjectService.DeleteProject",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return domainerror.NewProjectNotProvidedError(
			fmt.Errorf(ErrProjectIDNotProvided),
		)
	}

	project, err = s.repository.Find(projectID)
	if err != nil {
		s.logger.Error("%s: %s", ErrFindingProject, err.Error(), map[string]interface{}{
			"component":  "DeleteProjectService.DeleteProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": projectID,
		})
		return domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s: %w", ErrFindingProject, err),
		)
	}

	storer = s.storage.Get(project.Storage)
	if storer == nil {
		s.logger.Error(ErrStorageHandlerNotFound, map[string]interface{}{
			"component":  "DeleteProjectService.DeleteProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": projectID,
			"storage":    project.Storage,
		})
		return fmt.Errorf(ErrStorageHandlerNotFound)
	}

	err = s.repository.Delete(projectID)
	if err != nil {
		s.logger.Error("%s: %s", ErrDeletingProject, err.Error(), map[string]interface{}{
			"component":  "DeleteProjectService.DeleteProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": projectID,
		})
		return fmt.Errorf("%s: %w", ErrDeletingProject, err)
	}

	err = storer.Delete(project)
	if err != nil {
		s.logger.Error("%s: %s", ErrDeletingProject, err.Error(), map[string]interface{}{
			"component":  "DeleteProjectService.DeleteProject",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id": projectID,
			"storage":    project.Storage,
		})
		return fmt.Errorf("%s: %w", ErrDeletingProject, err)
	}

	return nil
}

// DeleteVersion deletes a project version by its id and version
func (s *DeleteProjectService) DeleteVersion(projectID string, version string) error {
	// TODO: implement the logic to delete a project version. This will require changes in the repository and storage layers to support versioning.
	return fmt.Errorf("DeleteVersion method not implemented yet")
}
