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
		errorMsg := fmt.Sprintf("error opening file %s: %s", source, err)
		s.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":   "LocalFetchFile.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})

		return fmt.Errorf("%s", errorMsg)
	}

	destPath := filepath.Join(workingDir, sourceFileInfo.Name())
	dstFile, err := s.fs.Create(destPath)
	defer func() {
		err = dstFile.Close()
	}()
	if err != nil {
		errorMsg := fmt.Sprintf("error creating file %s: %s", destPath, err)
		s.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":   "LocalFetchFile.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})
		return fmt.Errorf("%s", errorMsg)
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: From %s to %s: %s", ErrCopyingFilesInWorkingDir, source, destPath, err)
		s.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":   "LocalFetchFile.Fetch",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch",
				"source_dir":  source,
				"working_dir": workingDir,
			})
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
