package unpack

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/tar"
	"github.com/spf13/afero"
)

var (
	// ErrTarExtractorNotProvided is returned when the tar extractor is not provided
	ErrTarExtractorNotProvided = fmt.Errorf("Tar extractor not provided")
)

// TarGzipFormat struct used to unpack tar.gz files
type TarGzipFormat struct {
	fs        afero.Fs
	logger    repository.Logger
	extractor *tar.Tar
}

// NewTarGzipFormat method creates a new TarGzipFormat struct
func NewTarGzipFormat(fs afero.Fs, logger repository.Logger) *TarGzipFormat {
	return &TarGzipFormat{
		fs:        fs,
		logger:    logger,
		extractor: tar.NewTar(fs, logger),
	}
}

// Unpack method prepares the project into dest folder
func (a *TarGzipFormat) Unpack(project *entity.Project, workingDir string) error {
	var err error
	var gzipReader *gzip.Reader
	var sourceCodeFileReader io.Reader
	var sourceCodeFile string
	var sourceFileInfo os.FileInfo

	if project == nil {
		a.logger.Error(ErrProjectNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrProjectNotProvided
	}

	if a.extractor == nil {
		a.logger.Error(ErrTarExtractorNotProvided.Error(),
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
			})
		return ErrTarExtractorNotProvided
	}

	// It should be received instead of being figured out here
	sourceFileInfo, _ = a.fs.Stat(project.Reference)
	sourceCodeFile = filepath.Join(workingDir, sourceFileInfo.Name())

	_, err = a.fs.Stat(sourceCodeFile)
	// error if the file does not exists
	if err != nil {
		errorMsg := fmt.Sprintf("source code file %s does not exists", sourceCodeFile)
		a.logger.Error(errorMsg,
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				"file":      sourceCodeFile,
			})

		return fmt.Errorf("%s", errorMsg)
	}

	sourceCodeFileReader, err = a.fs.Open(sourceCodeFile)
	if err != nil {
		errorMsg := fmt.Sprintf("error opening source code file %s", sourceCodeFile)
		a.logger.Error(errorMsg,
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				"file":      sourceCodeFile,
			})
		return fmt.Errorf("%s", errorMsg)
	}

	gzipReader, err = gzip.NewReader(sourceCodeFileReader)
	if err != nil {
		errorMsg := fmt.Sprintf("error creating gzip reader from %s", sourceCodeFile)
		a.logger.Error(errorMsg,
			map[string]interface{}{
				"component": "TarGzipFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				"file":      sourceCodeFile,
			})
		return fmt.Errorf("%s", errorMsg)
	}

	err = a.extractor.Extract(gzipReader, workingDir)
	if err != nil {
		errorMsg := fmt.Sprintf("error extracting tar file %s into %s", sourceCodeFile, workingDir)
		a.logger.Error(errorMsg,
			map[string]interface{}{
				"component":   "TarGzipFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/archive",
				"file":        sourceCodeFile,
				"working_dir": workingDir,
			})
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
