package archive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

type LocalStorageArchive struct {
	Fs     afero.Fs
	logger repository.Logger
}

func NewLocalStorageArchive(fs afero.Fs, logger repository.Logger) *LocalStorageArchive {
	return &LocalStorageArchive{
		Fs:     fs,
		logger: logger,
	}
}

// Unarchive method prepares the project into dest folder
func (a *LocalStorageArchive) Unarchive(project *entity.Project, workingDir string) error {
	var err error

	if project == nil {
		a.logger.Error(
			ErrProjectNotProvided.Error(),
			map[string]interface{}{
				"component": "LocalStorageArchive.Unarchive",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrProjectNotProvided
	}

	_, err = a.Fs.Stat(workingDir)
	if err == nil {
		errorMsg := fmt.Sprintf("working directory %s already exists", workingDir)
		a.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "LocalStorageArchive.Unarchive",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return fmt.Errorf("%s", errorMsg)
	}

	// Create working directory
	err = a.Fs.MkdirAll(workingDir, 0755)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCreatingWorkingDirFolder, err)
		a.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "LocalStorageArchive.Unarchive",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return fmt.Errorf("%s", errorMsg)
	}

	a.logger.Info(fmt.Sprintf("Setup project %s to %s", project.Name, workingDir))

	err = afero.Walk(a.Fs, project.Reference, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			errorMsg := fmt.Sprintf("error walking through %s: %s", path, err)
			a.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorageArchive.Unarchive",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				})
			return fmt.Errorf("%s", errorMsg)
		}

		relPath, err := filepath.Rel(project.Reference, path)
		if err != nil {
			errorMsg := fmt.Sprintf("error getting relative path for %s: %s", path, err)
			a.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorageArchive.Unarchive",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				})
			return fmt.Errorf("%s", errorMsg)
		}

		// Create directories
		if info.IsDir() {
			err = a.Fs.MkdirAll(filepath.Join(workingDir, relPath), 0755)
			if err != nil {
				errorMsg := fmt.Sprintf("error creating directory %s: %s", filepath.Join(workingDir, relPath), err)
				a.logger.Error(
					errorMsg,
					map[string]interface{}{
						"component": "LocalStorageArchive.Unarchive",
						"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
					})
				return fmt.Errorf("%s", errorMsg)
			} else {
				return nil
			}
		}

		// Copy files
		srcFile, err := a.Fs.Open(path)
		if err != nil {
			errorMsg := fmt.Sprintf("error opening file %s: %s", path, err)
			a.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorageArchive.Unarchive",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				})
			return fmt.Errorf("%s", errorMsg)
		}
		defer srcFile.Close()

		destPath := filepath.Join(workingDir, relPath)
		destFile, err := a.Fs.Create(destPath)
		if err != nil {
			errorMsg := fmt.Sprintf("error creating file %s: %s", destPath, err)
			a.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorageArchive.Unarchive",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				})
			return fmt.Errorf("%s", errorMsg)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			errorMsg := fmt.Sprintf("%s: From %s to %s: %s", ErrCopyingFilesInWorkingDir, path, destPath, err)
			a.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component": "LocalStorageArchive.Unarchive",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				})
			return fmt.Errorf("%s", errorMsg)
		}

		return nil

	})

	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCopyingFilesInWorkingDir, err)
		a.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "LocalStorageArchive.Unarchive",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
