package error

// TaskNotProvidedError is an error type launched when a task is not provided
type TaskNotProvidedError struct {
	Err error
}

// NewTaskNotProvidedError creates a new TaskNotProvidedError
func NewTaskNotProvidedError(err error) *TaskNotProvidedError {
	return &TaskNotProvidedError{Err: err}
}

// Error returns the error message
func (e *TaskNotProvidedError) Error() string {
	return e.Err.Error()
}
