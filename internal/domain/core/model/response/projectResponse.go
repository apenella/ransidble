package response

// ProjectResponse represents a response describing a project
type ProjectResponse struct {
	// Format represents the project format
	Format string `json:"format" validate:"required"`
	// Name represents the project name
	Name string `json:"name" validate:"required"`
	// Source represents the project source
	Reference string `json:"reference" validate:"required"`
	// Storage represents the project type
	Storage string `json:"storage" validate:"required"`
}
