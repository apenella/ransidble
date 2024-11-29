package repository

import (
	"os"
)

// Workspacer interface to manage a workspace
type Workspacer interface {
	// Prepare prepares the workspace
	Prepare() error
	// GetWorkingDir returns the working directory
	GetWorkingDir() (string, error)
	// Cleanup cleans up the workspace
	Cleanup() error
}

// Filesystemer interface to manage a filesystem
type Filesystemer interface {
	// MkdirAll creates a directory
	MkdirAll(path string, perm os.FileMode) error
	// RemoveAll removes a directory
	RemoveAll(path string) error
	// Stat returns the file information
	Stat(path string) (os.FileInfo, error)
	// TempDir creates a temporary directory
	TempDir(dir, prefix string) (name string, err error)
}
