package archive

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

var (
	// ErrCreatingWorkingDirFolder represents an error creating working directory folder
	ErrCreatingWorkingDirFolder = errors.New("error creating working directory folder")
	// ErrCopyingFilesInWorkingDir represents an error copying files in working directory
	ErrCopyingFilesInWorkingDir = errors.New("error copying files in working directory")
	// ErrRemovingWorkingDirFolder represents an error removing working directory folder
	ErrRemovingWorkingDirFolder = errors.New("error removing working directory folder")
	// ErrProjectNotProvided represents an error when the project is not provided
	ErrProjectNotProvided = errors.New("project not provided")
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

// Setup method prepares the project into dst folder
func (a *LocalStorageArchive) Unarchive(project *entity.Project, workingDir string) error {
	var err error

	if project == nil {
		a.logger.Error(ErrProjectNotProvided.Error())
		return ErrProjectNotProvided
	}

	_, err = a.Fs.Stat(workingDir)
	if err == nil {
		errorMsg := fmt.Sprintf("working directory %s already exists", workingDir)
		a.logger.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Create working directory
	err = a.Fs.MkdirAll(workingDir, 0755)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCreatingWorkingDirFolder, err)
		a.logger.Error(errorMsg, "LocalStorageArchive::Unarchive")
		return fmt.Errorf(errorMsg)
	}

	a.logger.Info(fmt.Sprintf("Setup project %s to %s", project.Name, workingDir))

	err = afero.Walk(a.Fs, project.Reference, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			errorMsg := fmt.Sprintf("error walking through %s: %s", path, err)
			a.logger.Error(errorMsg)
			return fmt.Errorf(errorMsg)
		}

		relPath, err := filepath.Rel(project.Reference, path)
		if err != nil {
			errorMsg := fmt.Sprintf("error getting relative path for %s: %s", path, err)
			a.logger.Error(errorMsg)
			return fmt.Errorf(errorMsg)
		}

		// Create directories
		if info.IsDir() {
			err = a.Fs.MkdirAll(filepath.Join(workingDir, relPath), 0755)
			if err != nil {
				errorMsg := fmt.Sprintf("error creating directory %s: %s", filepath.Join(workingDir, relPath), err)
				a.logger.Error(errorMsg)
				return fmt.Errorf(errorMsg)
			} else {
				return nil
			}
		}

		// Copy files
		srcFile, err := a.Fs.Open(path)
		if err != nil {
			errorMsg := fmt.Sprintf("error opening file %s: %s", path, err)
			a.logger.Error(errorMsg)
			return fmt.Errorf(errorMsg)
		}
		defer srcFile.Close()

		destPath := filepath.Join(workingDir, relPath)
		destFile, err := a.Fs.Create(destPath)
		if err != nil {
			errorMsg := fmt.Sprintf("error creating file %s: %s", destPath, err)
			a.logger.Error(errorMsg)
			return fmt.Errorf(errorMsg)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			errorMsg := fmt.Sprintf("%s: From %s to %s: %s", ErrCopyingFilesInWorkingDir, path, destPath, err)
			a.logger.Error(errorMsg)
			return fmt.Errorf(errorMsg)
		}

		return nil

	})

	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCopyingFilesInWorkingDir, err)
		a.logger.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	return nil
}
