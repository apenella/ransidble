package error

// ProjectNotProvidedError is an error type launched when a project is not provided
type ProjectNotProvidedError struct {
	Err error
}

// NewProjectNotProvidedError creates a new ProjectNotProvidedError
func NewProjectNotProvidedError(err error) *ProjectNotProvidedError {
	return &ProjectNotProvidedError{Err: err}
}

// Error returns the error message
func (e *ProjectNotProvidedError) Error() string {
	return e.Err.Error()
}
