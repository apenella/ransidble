package response

// CommandResponse represents a response
type CommandResponse struct {
	ID string `json:"id" validate:"required"`
	// Status string `json:"status" validate:"required"`
	// Data   interface{} `json:"data" validate:"required"`
	Error string `json:"error,omitempty" validate:"string"`
}
