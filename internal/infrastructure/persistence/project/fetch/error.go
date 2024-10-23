package fetch

import "errors"

var (
	// ErrSourceCodeNotExists represents an error when the source code does not exists
	ErrSourceCodeNotExists = errors.New("source code does not exists")
	// ErrWorkingDirNotExists represents an error when the destination to fetch does not exists
	ErrWorkingDirNotExists = errors.New("working directory does not exists")
	// ErrCopyingFilesInWorkingDir represents an error copying files in working directory
	ErrCopyingFilesInWorkingDir = errors.New("error copying files in working directory")
)
