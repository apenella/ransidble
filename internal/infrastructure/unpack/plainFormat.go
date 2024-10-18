package unpack

import (
	"fmt"

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

var (
	// ErrWorkingDirNotExists represents an error when the destination to fetch does not exists
	ErrWorkingDirNotExists = fmt.Errorf("working directory does not exists")
)

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

	_, err = p.fs.Stat(workingDir)
	if err != nil {
		workingDirExist, err = afero.DirExists(p.fs, workingDir)
		if workingDirExist == false || err != nil {
			errorMsg := fmt.Errorf("%s. %w", ErrWorkingDirNotExists, err)
			p.logger.Error(
				errorMsg.Error(),
				map[string]interface{}{
					"component": "PlainFormat.Unpack",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/unpack",
				})

			return errorMsg
		}
	}

	return nil
}
