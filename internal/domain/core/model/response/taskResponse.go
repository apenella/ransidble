package response

// TaskResponse represents a response
type TaskResponse struct {
	// ID of the task created
	ID string `json:"id" validate:"required"`
	// Error represents an error
	Error string `json:"error,omitempty" validate:"string"`
}
