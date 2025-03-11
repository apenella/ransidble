package error

// ProjectAlreadyExistsError is an error type for project already exists
type ProjectAlreadyExistsError struct {
	Err error
}

// NewProjectAlreadyExistsError creates a new ProjectAlreadyExistsError
func NewProjectAlreadyExistsError(err error) *ProjectAlreadyExistsError {
	return &ProjectAlreadyExistsError{Err: err}
}

// Error returns the error message
func (e *ProjectAlreadyExistsError) Error() string {
	return e.Err.Error()
}
