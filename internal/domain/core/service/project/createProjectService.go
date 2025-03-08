package project

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

const (
	// ErrProjectFormatNotProvided error message when format is not provided
	ErrProjectFormatNotProvided = "format not provided"
	// ErrProjectStorageNotProvided error message when storage is not provided
	ErrProjectStorageNotProvided = "storage not provided"
	// ErrNotFoundStorageHandler error message when storage handler is not found
	ErrNotFoundStorageHandler = "storage handler not found"
	// ErrProjectFileNotProvided error message when file is not provided
	ErrProjectFileNotProvided = "file not provided"
	// ErrProjectStorageNotSupported error message when storage is not supported
	ErrProjectStorageNotSupported = "storage not supported"
	// ErrProjectFormatNotSupported error message when format is not supported
	ErrProjectFormatNotSupported = "format not supported"
	// ErrStoringProject error message when storing project fails
	ErrStoringProject = "storing project fails"
	// ErrProjectFileExtensionNotSupported error message when file extension is not supported
	ErrProjectFileExtensionNotSupported = "file extension not supported"
	// ErrStorageHandlerNotInitialized error message when storage handler is not initialized
	ErrStorageHandlerNotInitialized = "storage handler not initialized"
	// ErrProjectRepositoryNotInitialized error message when project repository is not initialized
	ErrProjectRepositoryNotInitialized = "project repository not initialized"
	// ErrOpeningProjectFile error message when opening project file fails
	ErrOpeningProjectFile = "opening project file fails"
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
			"file":      reference,
		})
		return fmt.Errorf("%s: %s", ErrProjectFileExtensionNotSupported, err.Error())
	}

	err = entity.ValidateProjectStorage(storage)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectStorageNotSupported, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"storage":   storage,
		})
		return fmt.Errorf("%s: %s", ErrProjectStorageNotSupported, err.Error())
	}

	err = entity.ValidateProjectFormat(format)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrProjectFormatNotSupported, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"format":    format,
		})
		return fmt.Errorf("%s: %s", ErrProjectFormatNotSupported, err.Error())
	}

	storer := s.storage.Get(storage)
	if storer == nil {
		s.logger.Error(ErrNotFoundStorageHandler, map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"storage":   storage,
		})
		return fmt.Errorf(ErrNotFoundStorageHandler)
	}

	name := extractProjectName(reference)

	projectFileReader, err = file.Open()
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrOpeningProjectFile, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
			"reference": reference,
		})
		return fmt.Errorf("%s: %s", ErrOpeningProjectFile, err.Error())
	}

	project := entity.NewProject(name, reference, format, storage)
	err = s.repository.SafeStore(name, project)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrStoringProject, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())
	}

	err = storer.Store(project, projectFileReader)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", ErrStoringProject, err.Error()), map[string]interface{}{
			"component": "CreateProjectService.Create",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/project",
		})
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())
	}

	return nil
}

// extractProjectName extracts the project name from the file
func extractProjectName(file string) string {
	return strings.Split(file, ".")[0]
}
