package fetch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

// LocalFetchFile represents the local fetcher for fetching source code components from a file
type LocalFetchFile struct {
	// fs is the filesystem
	fs afero.Fs
	// logger is the logger
	logger repository.Logger
}

// NewLocalFetchFile creates a new local fetch file
func NewLocalFetchFile(fs afero.Fs, logger repository.Logger) *LocalFetchFile {
	return &LocalFetchFile{
		fs:     fs,
		logger: logger,
	}
}

// Fetch method copies the project from local storage to working directory
func (s *LocalFetchFile) Fetch(source string, workingDir string) (err error) {

	var workingDirExist bool
	var sourceFileExist bool
	var sourceFileInfo os.FileInfo

	if s.fs == nil {
		s.logger.Error(
			ErrFileSystemNotInitialized.Error(),
			map[string]interface{}{
				"component": "LocalFetchFile.Fetch",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
			})
		return ErrFileSystemNotInitialized
	}

	sourceFileInfo, err = s.fs.Stat(source)
	if err != nil {
		sourceFileExist, err = afero.Exists(s.fs, source)
		if sourceFileExist == false || err != nil {
			s.logger.Error(
				ErrSourceCodeNotExists.Error(),
				map[string]interface{}{
					"component":  "LocalFetchFile.Fetch",
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
					"component":   "LocalFetchFile.Fetch",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
					"working_dir": workingDir,
				})
			return ErrWorkingDirNotExists
		}
	}

	srcFile, err := s.fs.Open(source)
	defer func() {
		err = srcFile.Close()
	}()
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrOpeningASourceCodeFileFromLocalDir, err),
			map[string]interface{}{
				"component":   "LocalFetchFile.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})

		return fmt.Errorf("%s: %w", ErrOpeningASourceCodeFileFromLocalDir, err)
	}

	destPath := filepath.Join(workingDir, sourceFileInfo.Name())
	dstFile, err := s.fs.Create(destPath)
	defer func() {
		err = dstFile.Close()
	}()

	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrCreatingAFileFromLocalToDirWorkingDir, err),
			map[string]interface{}{
				"component":   "LocalFetchFile.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})
		return fmt.Errorf("%s: %w", ErrCreatingAFileFromLocalToDirWorkingDir, err)
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("%s: %s", ErrCopyingAFileFromLocalToDirWorkingDir, err),
			map[string]interface{}{
				"component":        "LocalFetchFile.Fetch",
				"destination_path": destPath,
				"package":          "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_file":      source,
				"working_dir":      workingDir,
			})
		return fmt.Errorf("%s: %w", ErrCopyingAFileFromLocalToDirWorkingDir, err)
	}

	return nil
}
