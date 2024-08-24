package response

// TaskResponse represents a response describing a task
type TaskResponse struct {
	// Command identifies the command to be executed
	Command string `json:"command" validate:"required"`
	// CompletedAt represents the time the task was completed
	CompletedAt string `json:"completed_at"`
	// CreatedAt represents the time the task was created
	CreatedAt string `json:"created_at"`
	// ErrorMessage represents an error message
	ErrorMessage string `json:"error_message,omitempty"`
	// ExecutedAt represents the time the task was executed
	ExecutedAt string `json:"executed_at"`
	// ID represents the task ID
	ID string `json:"id" validate:"required"`
	// Parameters represents the parameters to be used
	Parameters interface{} `json:"parameters" validate:"required"`
	// Project represents the project
	Project *ProjectResponse `json:"project"`
	// Status represents the status of the task
	Status string `json:"status" validate:"required"`
}
