package response

// ProjectResponse represents a response describing a project
type ProjectResponse struct {
	// Name represents the project name
	Name string `json:"name" validate:"required"`
	// Source represents the project source
	Reference string `json:"reference" validate:"required"`
	// Type represents the project type
	Type string `json:"type" validate:"required"`
}
