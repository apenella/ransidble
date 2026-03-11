package project

import (
	"fmt"
	"io"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
)

// CreateProjectService represents the service to create a project
type CreateProjectService struct {
	repository repository.ProjectRepository
	storage    repository.SourceCodeStorageFactory
	logger     repository.Logger
}

// Ensure CreateProjectService implements the CreateProjectServicer interface
var _ service.CreateProjectServicer = (*CreateProjectService)(nil)

// NewCreateProjectService creates a new CreateProjectService
func NewCreateProjectService(repository repository.ProjectRepository, storage repository.SourceCodeStorageFactory, logger repository.Logger) *CreateProjectService {
	return &CreateProjectService{
		storage:    storage,
		repository: repository,
		logger:     logger,
	}
}

// Create creates a project and returns an error if something goes wrong
// func (s *CreateProjectService) Create(format string, storage string, file *multipart.FileHeader) error {
func (s *CreateProjectService) Create(format string, storage string, projectID string, projectVersion string, projectContentReader io.Reader) error {
	var err error
	var extension string

	if format == "" {
		s.logger.Error(ErrProjectFormatNotProvided, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf(ErrProjectFormatNotProvided)
	}

	if storage == "" {
		s.logger.Error(ErrProjectStorageNotProvided, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf(ErrProjectStorageNotProvided)
	}

	if projectContentReader == nil {
		s.logger.Error(ErrProjectContentReaderNotProvided, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf(ErrProjectContentReaderNotProvided)
	}

	if projectID == "" {
		s.logger.Error(fmt.Sprintf(ErrProjectIDNotProvided), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return domainerror.NewProjectIDNotProvidedError(
			fmt.Errorf(ErrProjectIDNotProvided),
		)
	}

	if s.storage == nil {
		s.logger.Error(ErrStorageHandlerNotInitialized, map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
		})
		return fmt.Errorf(ErrStorageHandlerNotInitialized)
	}

	if s.repository == nil {
		s.logger.Error(ErrProjectRepositoryNotInitialized, map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
		})
		return fmt.Errorf(ErrProjectRepositoryNotInitialized)
	}

	err = entity.ValidateProjectFormat(format)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectFormatNotSupported, err.Error()), map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"format":          format,
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
			"storage":         storage,
		})
		return fmt.Errorf("%s: %s", ErrProjectFormatNotSupported, err.Error())
	}

	err = entity.ValidateProjectStorage(storage)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectStorageNotSupported, err.Error()), map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
			"storage":         storage,
		})
		return fmt.Errorf("%s: %s", ErrProjectStorageNotSupported, err.Error())
	}

	extension, err = entity.GetExtensionFromFormat(format)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectFormatNotSupported, err.Error()), map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"format":          format,
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
			"storage":         storage,
		})
		return fmt.Errorf("%s: %s", ErrProjectFormatNotSupported, err.Error())
	}

	// name := extractProjectName(filename)
	findProject, _ := s.repository.Find(projectID)
	if findProject != nil {
		s.logger.Error(fmt.Sprintf(ErrProjectAlreadyExists), map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":          format,
			"project_id":      projectID,
			"project_version": projectVersion,
			"storage":         storage,
		})
		return domainerror.NewProjectAlreadyExistsError(
			fmt.Errorf(ErrProjectAlreadyExists),
		)
	}

	storer := s.storage.Get(storage)
	if storer == nil {
		s.logger.Error(ErrStorageHandlerNotFound, map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"format":          format,
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
			"storage":         storage,
		})
		return fmt.Errorf(ErrStorageHandlerNotFound)
	}

	reference := fmt.Sprintf("%s.%s", projectID, extension)

	project := entity.NewProject(projectID, projectVersion, reference, format, storage)
	err = s.repository.SafeStore(projectID, project)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrStoringProject, err.Error()), map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"format":          format,
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
			"reference":       reference,
			"storage":         storage,
		})
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())
	}

	err = storer.Store(project, projectContentReader)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrStoringProject, err.Error()), map[string]interface{}{
			"component":       "CreateProjectService.Create",
			"format":          format,
			"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
			"project_id":      projectID,
			"project_version": projectVersion,
			"reference":       reference,
			"storage":         storage,
		})
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())

	}

	s.logger.Info("Project created", map[string]interface{}{
		"component":       "CreateProjectService.Create",
		"package":         "github.com/apenella/ransidble/internal/domain/core/service/project",
		"format":          format,
		"project_id":      projectID,
		"project_version": projectVersion,
		"reference":       reference,
		"storage":         storage,
	})

	return nil
}
