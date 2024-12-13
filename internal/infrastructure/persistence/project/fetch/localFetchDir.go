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
func (l *LocalFetchDir) Fetch(source string, workingDir string) (err error) {

	var workingDirExist bool
	var sourceDirExist bool

	if l.fs == nil {
		l.logger.Error(
			ErrFileSystemNotInitialized.Error(),
			map[string]interface{}{
				"component": "LocalFetchDir.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})
		return ErrFileSystemNotInitialized
	}

	_, err = l.fs.Stat(source)
	if err != nil {
		sourceDirExist, err = afero.DirExists(l.fs, source)
		if sourceDirExist == false || err != nil {
			l.logger.Error(
				ErrSourceCodeNotExists.Error(),
				map[string]interface{}{
					"component":  "LocalFetchDir.Fetch",
					"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir": source,
				})
			return ErrSourceCodeNotExists
		}
	}

	_, err = l.fs.Stat(workingDir)
	if err != nil {
		workingDirExist, err = afero.DirExists(l.fs, workingDir)
		if workingDirExist == false || err != nil {
			l.logger.Error(
				ErrWorkingDirNotExists.Error(),
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"working_dir": workingDir,
				})
			return ErrWorkingDirNotExists
		}
	}

	err = afero.Walk(l.fs, source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			l.logger.Error(
				fmt.Sprintf("%s: %s", ErrWalkingDirToFetchSourceCodeFromLocalDir, err),
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"path":        path,
					"source_dir":  source,
					"working_dir": workingDir,
				})
			return fmt.Errorf("%s: %w", ErrWalkingDirToFetchSourceCodeFromLocalDir, err)
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			l.logger.Error(
				fmt.Sprintf("%s: %s", ErrGettingSourceCodeRelativePathFromLocalDir, err),
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"path":        path,
					"source_dir":  source,
					"working_dir": workingDir,
				})
			return fmt.Errorf("%s: %w", ErrGettingSourceCodeRelativePathFromLocalDir, err)
		}

		if source == path && info.IsDir() {
			return nil
		}

		if info.IsDir() {
			err = l.fs.MkdirAll(filepath.Join(workingDir, relPath), 0755)
			if err != nil {
				l.logger.Error(
					fmt.Sprintf("%s: %s", ErrCreatingADirectoryInLocalDirWorkingDir, err),
					map[string]interface{}{
						"component":     "LocalFetchDir.Fetch",
						"package":       "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
						"path":          path,
						"relative_path": relPath,
						"source_dir":    source,
						"working_dir":   workingDir,
					})

				return fmt.Errorf("%s: %w", ErrCreatingADirectoryInLocalDirWorkingDir, err)

			}
		}

		srcFile, err := l.fs.Open(path)
		defer func() {
			err = srcFile.Close()
		}()
		if err != nil {
			// errorMsg := fmt.Sprintf("error opening file %s: %s", path, err)
			l.logger.Error(
				fmt.Sprintf("%s: %s", ErrOpeningASourceCodeFileFromLocalDir, err),
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"path":        path,
					"source_dir":  source,
					"working_dir": workingDir,
				})

			return fmt.Errorf("%s: %w", ErrOpeningASourceCodeFileFromLocalDir, err)
		}

		destPath := filepath.Join(workingDir, relPath)
		destFile, err := l.fs.Create(destPath)
		defer func() {
			err = destFile.Close()
		}()
		if err != nil {
			l.logger.Error(
				fmt.Sprintf("%s: %s", ErrCreatingAFileFromLocalToDirWorkingDir, err),
				map[string]interface{}{
					"component":   "LocalFetchDir.Fetch",
					"destiantion": destPath,
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"path":        path,
					"source_dir":  source,
					"working_dir": workingDir,
				})

			return fmt.Errorf("%s: %w", ErrCreatingAFileFromLocalToDirWorkingDir, err)
		}

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			l.logger.Error(
				fmt.Sprintf("%s: %s", ErrCopyingAFileFromLocalToDirWorkingDir, err),
				map[string]interface{}{
					"component":        "LocalFetchDir.Fetch",
					"destiantion_path": destPath,
					"package":          "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"source_dir":       source,
					"source_file":      path,
					"working_dir":      workingDir,
				})

			return fmt.Errorf("%s: %w", ErrCopyingAFileFromLocalToDirWorkingDir, err)
		}

		return nil

	})

	if err != nil {
		l.logger.Error(
			fmt.Sprintf("%s: %s", ErrCopyingFilesToWorkingDir, err),
			map[string]interface{}{
				"component":   "LocalFetchDir.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})

		return fmt.Errorf("%s: %w", ErrCopyingFilesToWorkingDir, err)
	}

	return nil
}
