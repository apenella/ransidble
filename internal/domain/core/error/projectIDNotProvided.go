package error

// ProjectIDNotProvidedError is an error type for project id not provided
type ProjectIDNotProvidedError struct {
	Err error
}

// NewProjectIDNotProvidedError creates a new ProjectIDNotProvidedError
func NewProjectIDNotProvidedError(err error) *ProjectIDNotProvidedError {
	return &ProjectIDNotProvidedError{Err: err}
}

// Error returns the error message
func (e *ProjectIDNotProvidedError) Error() string {
	return e.Err.Error()
}
