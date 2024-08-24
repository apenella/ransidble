package response

// TaskCreatedResponse represents a response when a task is created
type TaskCreatedResponse struct {
	// ID of the task created
	ID string `json:"id" validate:"required"`
}
