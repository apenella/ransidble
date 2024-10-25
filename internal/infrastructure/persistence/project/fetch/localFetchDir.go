package fetch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

// LocalFetchDir struct for fetching project from local directory
type LocalFetchDir struct {
	// fs is the filesystem
	fs afero.Fs
	// logger is the logger
	logger repository.Logger
}

// NewLocalFetchDir creates a new local fetch directory
func NewLocalFetchDir(fs afero.Fs, logger repository.Logger) *LocalFetchDir {
	return &LocalFetchDir{
		fs:     fs,
		logger: logger,
	}
}

// Fetch method copies the project from local storage to working directory
func (s *LocalFetchDir) Fetch(source string, workingDir string) (err error) {

	var workingDirExist bool
	var sourceDirExist bool

	_, err = s.fs.Stat(source)
	if err != nil {
		sourceDirExist, err = afero.DirExists(s.fs, source)
		if sourceDirExist == false || err != nil {
			s.logger.Error(
				ErrSourceCodeNotExists.Error(),
				map[string]interface{}{
					"component":  "LocalFetchDir.Fetch",
					"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir": source,
				})
			return ErrSourceCodeNotExists
		}
	}

	_, err = s.fs.Stat(workingDir)
	if err != nil {
		workingDirExist, err = afero.DirExists(s.fs, workingDir)
		if workingDirExist == false || err != nil {
			s.logger.Error(
				ErrWorkingDirNotExists.Error(),
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"working_dir": workingDir,
				})
			return ErrWorkingDirNotExists
		}
	}

	err = afero.Walk(s.fs, source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errorMsg := fmt.Sprintf("error walking through %s: %s", path, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir":  source,
					"working_dir": workingDir,
				})
			return fmt.Errorf("%s", errorMsg)
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			errorMsg := fmt.Sprintf("error getting relative path for %s: %s", path, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir":  source,
					"working_dir": workingDir,
				})
			return fmt.Errorf("%s", errorMsg)
		}

		if source == path && info.IsDir() {
			return nil
		}

		if info.IsDir() {
			err = s.fs.MkdirAll(filepath.Join(workingDir, relPath), 0755)
			if err != nil {
				errorMsg := fmt.Sprintf("error creating directory %s: %s", filepath.Join(workingDir, relPath), err)
				s.logger.Error(
					errorMsg,
					map[string]interface{}{
						"component":   "LocalFetchDir.Fetch",
						"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
						"source_dir":  source,
						"working_dir": workingDir,
					})

				return fmt.Errorf("%s", errorMsg)

			}
		}

		srcFile, err := s.fs.Open(path)
		defer func() {
			err = srcFile.Close()
		}()
		if err != nil {
			errorMsg := fmt.Sprintf("error opening file %s: %s", path, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir":  source,
					"working_dir": workingDir,
				})

			return fmt.Errorf("%s", errorMsg)
		}

		destPath := filepath.Join(workingDir, relPath)
		destFile, err := s.fs.Create(destPath)
		defer func() {
			err = destFile.Close()
		}()
		if err != nil {
			errorMsg := fmt.Sprintf("error creating file %s: %s", destPath, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir":  source,
					"working_dir": workingDir,
				})

			return fmt.Errorf("%s", errorMsg)
		}

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			errorMsg := fmt.Sprintf("%s: From %s to %s: %s", ErrCopyingFilesInWorkingDir, path, destPath, err)
			s.logger.Error(
				errorMsg,
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir":  source,
					"working_dir": workingDir,
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
				"component":   "LocalFetchDir.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})

		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
