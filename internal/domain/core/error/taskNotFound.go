package error

// TaskNotFoundError is an error type for project not found
type TaskNotFoundError struct {
	Err error
}

// NewTaskNotFoundError creates a new TaskNotFoundError
func NewTaskNotFoundError(err error) *TaskNotFoundError {
	return &TaskNotFoundError{Err: err}
}

// Error returns the error message
func (e *TaskNotFoundError) Error() string {
	return e.Err.Error()
}
