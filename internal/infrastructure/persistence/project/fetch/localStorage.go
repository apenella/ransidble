package fetch

import (
	"fmt"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

// LocalStorage represents a repository on local storage
type LocalStorage struct {
	// Filesystem path where projects are stored
	fs afero.Fs
	// logger is the logger
	logger repository.Logger
	// path where projects are stored
	path string
}

// NewLocalStorage creates a new local project repository
func NewLocalStorage(fs afero.Fs, path string, logger repository.Logger) *LocalStorage {
	return &LocalStorage{
		fs:     fs,
		path:   path,
		logger: logger,
	}
}

// Fetch method copies the project from local storage to working directory
func (s *LocalStorage) Fetch(project *entity.Project, workingDir string) (err error) {

	var sourceCodeFetcher SourceCodeFetcher
	var workingDirExist bool

	if project == nil {
		s.logger.Error(
			ErrProjectNotProvided.Error(),
			map[string]interface{}{
				"component": "LocalStorage.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})
		return ErrProjectNotProvided
	}

	if workingDir == "" {
		s.logger.Error(
			ErrWorkingDirNotProvided.Error(),
			map[string]interface{}{
				"component": "LocalStorage.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})
		return ErrWorkingDirNotProvided
	}

	if s.fs == nil {
		s.logger.Error(
			ErrFileSystemNotInitialized.Error(),
			map[string]interface{}{
				"component": "LocalStorage.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})
		return ErrFileSystemNotInitialized
	}

	_, err = s.fs.Stat(workingDir)
	if err != nil {
		workingDirExist, err = afero.DirExists(s.fs, workingDir)
		if workingDirExist == false || err != nil {
			s.logger.Error(
				ErrWorkingDirNotExists.Error(),
				map[string]interface{}{
					"component":   "LocalStorage.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"working_dir": workingDir,
				})
			return ErrWorkingDirNotExists
		}
	}

	if project.Reference == "" {
		s.logger.Error(
			ErrProjectReferenceNotProvided.Error(),
			map[string]interface{}{
				"component":   "LocalStorage.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"project_id":  project.Name,
				"working_dir": workingDir,
			})

		return ErrProjectReferenceNotProvided
	}

	souceCodeStorageLocation := filepath.Join(s.path, project.Reference)
	infoProjectReference, err := s.fs.Stat(souceCodeStorageLocation)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrSourceCodeNotExists, err),
			err.Error(),
			map[string]interface{}{
				"component": "LocalStorage.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})
		return fmt.Errorf("%s: %w", ErrSourceCodeNotExists, err)
	}

	if infoProjectReference.IsDir() {
		sourceCodeFetcher = NewLocalFetchDir(s.fs, s.logger)
	} else {
		sourceCodeFetcher = NewLocalFetchFile(s.fs, s.logger)
	}

	s.logger.Debug("fetching project", map[string]interface{}{
		"component":   "LocalStorage.Fetch",
		"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
		"project_id":  project.Name,
		"working_dir": workingDir,
	})

	err = sourceCodeFetcher.Fetch(souceCodeStorageLocation, workingDir)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrFetchingProjectFromLocalStorage, err),
			map[string]interface{}{
				"component": "LocalStorage.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})

		return fmt.Errorf("%s: %w", ErrFetchingProjectFromLocalStorage, err)
	}

	return nil
}
