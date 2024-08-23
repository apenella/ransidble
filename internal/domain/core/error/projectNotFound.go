package error

// ProjectNotFoundError is an error type for project not found
type ProjectNotFoundError struct {
	Err error
}

// NewProjectNotFoundError creates a new ProjectNotFoundError
func NewProjectNotFoundError(err error) *ProjectNotFoundError {
	return &ProjectNotFoundError{Err: err}
}

// Error returns the error message
func (e *ProjectNotFoundError) Error() string {
	return e.Err.Error()
}
