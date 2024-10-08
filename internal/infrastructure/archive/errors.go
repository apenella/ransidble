package archive

import "errors"

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
