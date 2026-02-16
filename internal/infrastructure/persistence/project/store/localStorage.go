package store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

const (
	// ErrProjectNotProvided represents the error when a project is not provided
	ErrProjectNotProvided = "project not provided"
	// ErrProjectFileNotProvided represents the error when a file is not provided
	ErrProjectFileNotProvided = "file not provided"
	// ErrProjectReferenceNotProvided represents the error when a project reference is not provided
	ErrProjectReferenceNotProvided = "project reference not provided"
	// ErrOpeningDestinationFileInLocalStorage represents the error when a destination file cannot be opened
	ErrOpeningDestinationFileInLocalStorage = "error opening destination file in local storage"
	// ErrStoringProjectInLocalStorage represents the error when a project cannot be stored in local storage
	ErrStoringProjectInLocalStorage = "error storing project in local storage"
	// ErrStorageHandlerNotInitialized represents the error when the storage filesystem is not initialized
	ErrStorageHandlerNotInitialized = "storage handler not initialized"
	// ErrStoragePathNotProvided represents the error when the storage path is not provided
	ErrStoragePathNotProvided = "storage path not provided"
	// ErrStoragePathNotExists represents the error when the storage path does not exists
	ErrStoragePathNotExists = "storage path not exists"
	// ErrStoragePathNotDirectory represents the error when the storage path is not a directory
	ErrStoragePathNotDirectory = "storage path not a directory"
	// ErrCheckingStoragePathIsDirectory represents the error when the storage path is not a directory
	ErrCheckingStoragePathIsDirectory = "error checking storage path is a directory"
	// ErrInitializingLocalStorage represents the error when the local storage cannot be initialized
	ErrInitializingLocalStorage = "error initializing local storage"
)

// LocalStorage represents a repository on local storage
type LocalStorage struct {
	// Filesystem path where projects are stored
	fs afero.Fs
	// path where projects are stored
	path string
	// logger is the logger
	logger repository.Logger
}

// NewLocalStorage creates a new local project repository
func NewLocalStorage(fs afero.Fs, path string, logger repository.Logger) *LocalStorage {
	return &LocalStorage{
		fs:     fs,
		path:   path,
		logger: logger,
	}
}

// Initialize method initializes the local storage
func (s *LocalStorage) Initialize() error {

	if s.fs == nil {
		s.logger.Error(
			ErrStorageHandlerNotInitialized,
			map[string]interface{}{
				"component": "LocalStorage.Initialize",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrStorageHandlerNotInitialized)
	}

	if s.path == "" {
		s.logger.Error(
			ErrStoragePathNotProvided,
			map[string]interface{}{
				"component": "LocalStorage.Initialize",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrStoragePathNotProvided)
	}

	_, err := s.fs.Stat(s.path)
	if err != nil {
		s.logger.Info(
			"Creating local project storage path",
			map[string]interface{}{
				"component": "LocalStorage.Initialize",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
				"path":      s.path,
			})

		err = s.fs.MkdirAll(s.path, 0755)
		if err != nil {
			s.logger.Error(
				fmt.Sprintf("%s: %s", ErrInitializingLocalStorage, err.Error()),
				map[string]interface{}{
					"component": "LocalStorage.Initialize",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
					"path":      s.path,
				})
			return fmt.Errorf("%s: %w", ErrInitializingLocalStorage, err)
		}
	}

	return nil
}

// Store method copies the project from working directory to local storage
func (s *LocalStorage) Store(project *entity.Project, srcFile io.Reader) (err error) {

	var dstFile afero.File

	if project == nil {
		s.logger.Error(
			ErrProjectNotProvided,
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrProjectNotProvided)
	}

	if srcFile == nil {
		s.logger.Error(
			ErrProjectFileNotProvided,
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrProjectFileNotProvided)
	}

	if project.Reference == "" {
		s.logger.Error(
			ErrProjectReferenceNotProvided,
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrProjectReferenceNotProvided)
	}

	if s.fs == nil {
		s.logger.Error(
			ErrStorageHandlerNotInitialized,
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrStorageHandlerNotInitialized)
	}

	if s.path == "" {
		s.logger.Error(
			ErrStoragePathNotProvided,
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
			})
		return fmt.Errorf(ErrStoragePathNotProvided)
	}

	_, err = s.fs.Stat(s.path)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrStoragePathNotExists, err.Error()),
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
				"path":      s.path,
			})
		return fmt.Errorf("%s: %s", ErrStoragePathNotExists, err.Error())
	}

	pathIsDir, err := afero.IsDir(s.fs, s.path)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrCheckingStoragePathIsDirectory, err.Error()),
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
				"path":      s.path,
			})
		return fmt.Errorf("%s: %s", ErrCheckingStoragePathIsDirectory, err.Error())
	}

	if !pathIsDir {
		s.logger.Error(
			ErrStoragePathNotDirectory,
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
				"path":      s.path,
			})
		return fmt.Errorf(ErrStoragePathNotDirectory)
	}

	destFilePath := filepath.Join(s.path, project.Reference)

	dstFile, err = s.fs.OpenFile(destFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrOpeningDestinationFileInLocalStorage, err.Error()),
			map[string]interface{}{
				"component": "LocalStorage.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
				"reference": project.Reference,
			})
		return fmt.Errorf("%s: %s", ErrOpeningDestinationFileInLocalStorage, err.Error())
	}

	defer func() {
		err = dstFile.Close()
	}()

	err = afero.WriteReader(s.fs, destFilePath, srcFile)
	if err != nil {
		s.logger.Error(
			ErrStoringProjectInLocalStorage,
			map[string]interface{}{
				"component":   "LocalStorage.Store",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/store",
				"destination": destFilePath,
			})
		return fmt.Errorf(ErrStoringProjectInLocalStorage)
	}

	return nil
}
