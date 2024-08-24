package response

// TaskErrorResponse represents a response when a there is an error
type TaskErrorResponse struct {
	// ID of the task created
	ID string `json:"id" validate:"required"`
	// Error represents an error
	Error string `json:"error,omitempty" validate:"string"`
}
