package unpack

import (
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

// TarGzipFormat struct used to unpack tar.gz files
type TarGzipFormat struct {
	fs        afero.Fs
	logger    repository.Logger
	extractor TarExtractorer
}

// NewTarGzipFormat method creates a new TarGzipFormat struct
func NewTarGzipFormat(fs afero.Fs, extractor TarExtractorer, logger repository.Logger) *TarGzipFormat {
	return &TarGzipFormat{
		fs:        fs,
		logger:    logger,
		extractor: extractor,
	}
}

// Unpack method prepares the project into dest folder
func (a *TarGzipFormat) Unpack(project *entity.Project, workingDir string) error {
	var err error
	var gzipReader *gzip.Reader
	var sourceCodeFileReader io.Reader
	var sourceCodeFile string
	// var sourceFileInfo os.FileInfo

	if project == nil {
		a.logger.Error(ErrProjectNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrProjectNotProvided
	}

	if workingDir == "" {
		a.logger.Error(ErrWorkingDirNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrWorkingDirNotProvided
	}

	if a.fs == nil {
		a.logger.Error(ErrFilesystemNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrFilesystemNotProvided
	}

	if a.extractor == nil {
		a.logger.Error(ErrTarExtractorNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrTarExtractorNotProvided
	}

	if project.Reference == "" {
		a.logger.Error(ErrProjectReferenceNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrProjectReferenceNotProvided
	}

	sourceCodeFile = filepath.Join(workingDir, project.Reference)
	_, err = a.fs.Stat(sourceCodeFile)
	if err != nil {
		// error if the file does not exists
		a.logger.Error(
			ErrSourceCodeFileNotExist.Error(),
			map[string]interface{}{
				"component":   "TarGzipFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/archive",
				"source_file": sourceCodeFile,
			})

		return ErrSourceCodeFileNotExist
	}

	sourceCodeFileReader, err = a.fs.Open(sourceCodeFile)
	if err != nil {
		a.logger.Error(
			fmt.Sprintf("%s: %s", ErrOpeningSourceCodeFile, err),
			map[string]interface{}{
				"component":   "TarGzipFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/archive",
				"source_file": sourceCodeFile,
			})
		return fmt.Errorf("%s: %w", ErrOpeningSourceCodeFile, err)
	}

	gzipReader, err = gzip.NewReader(sourceCodeFileReader)
	if err != nil {
		a.logger.Error(
			fmt.Sprintf("%s: %s", ErrCreatingGzipReader, err),
			map[string]interface{}{
				"component":   "TarGzipFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/archive",
				"source_file": sourceCodeFile,
			})
		return fmt.Errorf("%s: %w", ErrCreatingGzipReader, err)
	}

	err = a.extractor.Extract(gzipReader, workingDir)
	if err != nil {
		a.logger.Error(
			fmt.Sprintf("%s: %s", ErrExtractingSourceCodeFile, err),
			map[string]interface{}{
				"component":   "TarGzipFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/archive",
				"source_file": sourceCodeFile,
				"working_dir": workingDir,
			})
		return fmt.Errorf("%s: %w", ErrExtractingSourceCodeFile, err)
	}

	return nil
}
