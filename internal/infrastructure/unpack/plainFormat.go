package unpack

import (
	"fmt"
	"os"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

// PlainFormat struct used to unpack local files
type PlainFormat struct {
	// fs is the filesystem
	fs afero.Fs
	// logger is the logger
	logger repository.Logger
}

// NewPlainFormat method creates a new PlainFormat struct
func NewPlainFormat(fs afero.Fs, logger repository.Logger) *PlainFormat {
	return &PlainFormat{
		fs:     fs,
		logger: logger,
	}
}

// Unpack method prepares the project into the working directory. Unpacking a plain format project does not require any action. The project is already fetched in the working directory. It just checks if the working directory exists.
func (p *PlainFormat) Unpack(project *entity.Project, workingDir string) error {

	var err error
	var workingDirExist bool

	if project == nil {
		p.logger.Error(
			ErrProjectNotProvided.Error(),
			map[string]interface{}{
				"component": "PlainFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/unpack",
			})
		return ErrProjectNotProvided
	}

	if workingDir == "" {
		p.logger.Error(
			ErrWorkingDirNotProvided.Error(),
			map[string]interface{}{
				"component": "PlainFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/unpack",
			})
		return ErrWorkingDirNotProvided
	}

	if p.fs == nil {
		p.logger.Error(
			ErrFilesystemNotProvided.Error(),
			map[string]interface{}{
				"component": "PlainFormat.Unpack",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/unpack",
			})
		return ErrFilesystemNotProvided
	}

	_, err = p.fs.Stat(workingDir)
	if err != nil {

		if os.IsNotExist(err) {
			p.logger.Error(
				fmt.Sprintf("%s: %s", ErrWorkingDirNotExists, err),
				map[string]interface{}{
					"component":   "PlainFormat.Unpack",
					"package":     "github.com/apenella/ransidble/internal/infrastructure/unpack",
					"working_dir": workingDir,
				})

			return fmt.Errorf("%s: %w", ErrWorkingDirNotExists, err)
		}

		p.logger.Error(
			fmt.Sprintf("%s: %s", ErrDescribingWorkingDir, err),
			map[string]interface{}{
				"component":   "PlainFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/unpack",
				"working_dir": workingDir,
			})

		return fmt.Errorf("%s: %w", ErrDescribingWorkingDir, err)
	}

	workingDirExist, _ = afero.IsDir(p.fs, workingDir)
	if workingDirExist == false {
		p.logger.Error(
			fmt.Sprintf("%s: %s", ErrWorkingDirIsNotDirectory, err),
			map[string]interface{}{
				"component":   "PlainFormat.Unpack",
				"package":     "github.com/apenella/ransidble/internal/infrastructure/unpack",
				"working_dir": workingDir,
			})

		return ErrWorkingDirIsNotDirectory
	}

	return nil
}
