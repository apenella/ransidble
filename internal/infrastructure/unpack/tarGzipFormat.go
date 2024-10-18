package unpack

import (
	"fmt"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

// TarGzipFormat struct used to unpack tar.gz files
type TarGzipFormat struct {
	Fs     afero.Fs
	logger repository.Logger
}

// NewTarGzipFormat method creates a new TarGzipFormat struct
func NewTarGzipFormat(fs afero.Fs, logger repository.Logger) *TarGzipFormat {
	return &TarGzipFormat{
		Fs:     fs,
		logger: logger,
	}
}

// Unpack method prepares the project into dest folder
func (a *TarGzipFormat) Unpack(project *entity.Project, workingDir string) error {
	var err error

	if project == nil {
		a.logger.Error(ErrProjectNotProvided.Error(),
			map[string]interface{}{
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				"component": "TarGzipFormat.Unpack"})
		return ErrProjectNotProvided
	}

	_, err = a.Fs.Stat(workingDir)
	if err == nil {
		errorMsg := fmt.Sprintf("working directory %s already exists", workingDir)
		a.logger.Error(errorMsg,
			map[string]interface{}{
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				"component": "TarGzipFormat.Unpack"})
		return fmt.Errorf("%s", errorMsg)
	}

	// untar directories, file, etc from the project specified at project.Reference into workingDir
	err = a.Fs.MkdirAll(workingDir, 0755)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCreatingWorkingDirFolder, err)
		a.logger.Error(
			errorMsg,
			map[string]interface{}{
				"package":   "github.com/apenella/ransidble/internal/infrastructure/archive",
				"component": "TarGzipFormat.Unpack"})
		return fmt.Errorf("%s", errorMsg)
	}

	a.logger.Info(fmt.Sprintf("Setup project %s to %s", project.Name, workingDir))

	return nil
}
