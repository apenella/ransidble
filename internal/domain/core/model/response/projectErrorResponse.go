package response

// ProjectErrorResponse represents a response when a there is an error
type ProjectErrorResponse struct {
	// ID of the task created
	ID string `json:"id" validate:"required"`
	// Error represents an error
	Error string `json:"error,omitempty" validate:"string"`
	// Status represents the status of the response
	Status int `json:"status" validate:"required,number"`
}
