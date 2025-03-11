package project

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

// CreateProjectService represents the service to create a project
type CreateProjectService struct {
	repository repository.ProjectRepository
	storage    repository.SourceCodeStorageFactory
	logger     repository.Logger
}

// NewCreateProjectService creates a new CreateProjectService
func NewCreateProjectService(repository repository.ProjectRepository, storage repository.SourceCodeStorageFactory, logger repository.Logger) *CreateProjectService {
	return &CreateProjectService{
		storage:    storage,
		repository: repository,
		logger:     logger,
	}
}

// Create creates a project and returns an error if something goes wrong
func (s *CreateProjectService) Create(format string, storage string, file *multipart.FileHeader) error {
	var err error
	var projectFileReader io.Reader

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

	if file == nil {
		s.logger.Error(ErrProjectFileNotProvided, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf(ErrProjectFileNotProvided)
	}

	if s.storage == nil {
		s.logger.Error(ErrStorageHandlerNotInitialized, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf(ErrStorageHandlerNotInitialized)
	}

	if s.repository == nil {
		s.logger.Error(ErrProjectRepositoryNotInitialized, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf(ErrProjectRepositoryNotInitialized)
	}

	reference := file.Filename
	err = entity.ValidateProjectFileExtension(reference)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectFileExtensionNotSupported, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"reference": reference,
		})
		return fmt.Errorf("%s: %s", ErrProjectFileExtensionNotSupported, err.Error())
	}

	name := extractProjectName(reference)
	findProject, _ := s.repository.Find(name)
	if findProject != nil {
		s.logger.Error(fmt.Sprintf(ErrProjectAlreadyExists), map[string]interface{}{
			"component":  "CreateProjectService.Create",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":     format,
			"project_id": name,
			"reference":  reference,
			"storage":    storage,
		})
		return domainerror.NewProjectAlreadyExistsError(
			fmt.Errorf(ErrProjectAlreadyExists),
		)
	}

	err = entity.ValidateProjectStorage(storage)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectStorageNotSupported, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"reference": reference,
			"storage":   storage,
		})
		return fmt.Errorf("%s: %s", ErrProjectStorageNotSupported, err.Error())
	}

	storer := s.storage.Get(storage)
	if storer == nil {
		s.logger.Error(ErrStorageHandlerNotFound, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":    format,
			"reference": reference,
			"storage":   storage,
		})
		return fmt.Errorf(ErrStorageHandlerNotFound)
	}

	err = entity.ValidateProjectFormat(format)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectFormatNotSupported, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":    format,
			"reference": reference,
			"storage":   storage,
		})
		return fmt.Errorf("%s: %s", ErrProjectFormatNotSupported, err.Error())
	}

	projectFileReader, err = file.Open()
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrOpeningProjectFile, err.Error()), map[string]interface{}{
			"component":  "CreateProjectService.Create",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":     format,
			"project_id": name,
			"reference":  reference,
			"storage":    storage,
		})
		return fmt.Errorf("%s: %s", ErrOpeningProjectFile, err.Error())
	}

	project := entity.NewProject(name, reference, format, storage)
	err = s.repository.SafeStore(name, project)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrStoringProject, err.Error()), map[string]interface{}{
			"component":  "CreateProjectService.Create",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":     format,
			"project_id": name,
			"reference":  reference,
			"storage":    storage,
		})
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())
	}

	err = storer.Store(project, projectFileReader)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrStoringProject, err.Error()), map[string]interface{}{
			"component":  "CreateProjectService.Create",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":     format,
			"project_id": name,
			"reference":  reference,
			"storage":    storage,
		})
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())
	}

	s.logger.Info("Project created", map[string]interface{}{
		"component":  "CreateProjectService.Create",
		"package":    "github.com/apenella/ransidble/internal/domain/core/service/project",
		"format":     format,
		"project_id": name,
		"reference":  reference,
		"storage":    storage,
	})

	return nil
}

// extractProjectName extracts the project name from the file
func extractProjectName(file string) string {
	return strings.Split(file, ".")[0]
}
