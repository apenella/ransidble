package fetch

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

const (
	// ExtensionTarGz represents the tar.gz extension
	ExtensionTarGz = ".tar.gz"
)

var (
	// ErrProjectNotProvided represents an error when the project is not provided
	ErrProjectNotProvided = errors.New("project not provided")
	// ErrCopyingFilesInWorkingDir represents an error copying files in working directory
	ErrCopyingFilesInWorkingDir = errors.New("error copying files in working directory")
	// ErrWorkingDirNotExists represents an error when the destination to fetch does not exists
	ErrWorkingDirNotExists = errors.New("working directory does not exists")
)

// LocalStorage represents a repository on local storage
type LocalStorage struct {
	// Filesystem path where projects are stored
	Fs afero.Fs

	// logger is the logger
	logger repository.Logger
}

// NewLocalStorage creates a new local project repository
func NewLocalStorage(fs afero.Fs, logger repository.Logger) *LocalStorage {
	return &LocalStorage{
		Fs:     fs,
		logger: logger,
	}
}

// Fetch method copies the project from local storage to working directory
func (s *LocalStorage) Fetch(project *entity.Project, workingDir string) error {

	var err error
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

	_, err = s.Fs.Stat(workingDir)
	if err != nil {
		workingDirExist, err = afero.DirExists(s.Fs, workingDir)
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

	s.logger.Info("fetching project", map[string]interface{}{
		"component":   "LocalStorage.Fetch",
		"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
		"project_id":  project.Name,
		"working_dir": workingDir,
	})

	err = afero.Walk(s.Fs, project.Reference, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			errorMsg := fmt.Sprintf("error walking through %s: %s", path, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorage.Fetch",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				})
			return fmt.Errorf("%s", errorMsg)
		}

		relPath, err := filepath.Rel(project.Reference, path)
		if err != nil {
			errorMsg := fmt.Sprintf("error getting relative path for %s: %s", path, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorage.Fetch",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				})
			return fmt.Errorf("%s", errorMsg)
		}

		if project.Reference == path {
			return nil
		}

		if info.IsDir() {
			err = s.Fs.MkdirAll(filepath.Join(workingDir, relPath), 0755)
			if err != nil {
				errorMsg := fmt.Sprintf("error creating directory %s: %s", filepath.Join(workingDir, relPath), err)
				s.logger.Error(
					errorMsg,
					map[string]interface{}{
						"component": "LocalStorage.Fetch",
						"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					})

				return fmt.Errorf("%s", errorMsg)

			}
		}

		srcFile, err := s.Fs.Open(path)
		if err != nil {
			errorMsg := fmt.Sprintf("error opening file %s: %s", path, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorage.Fetch",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				})

			return fmt.Errorf("%s", errorMsg)
		}
		defer srcFile.Close()

		destPath := filepath.Join(workingDir, relPath)
		destFile, err := s.Fs.Create(destPath)
		if err != nil {
			errorMsg := fmt.Sprintf("error creating file %s: %s", destPath, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorage.Fetch",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				})

			return fmt.Errorf("%s", errorMsg)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			errorMsg := fmt.Sprintf("%s: From %s to %s: %s", ErrCopyingFilesInWorkingDir, path, destPath, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorage.Fetch",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				})

			return fmt.Errorf("%s", errorMsg)
		}

		return nil

	})

	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCopyingFilesInWorkingDir, err)
		s.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "LocalStorage.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})

		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
