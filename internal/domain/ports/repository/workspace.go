package repository

// Workspacer interface to manage a workspace
type Workspacer interface {
	Prepare() error
	GetWorkingDir() (string, error)
	Cleanup() error
}
