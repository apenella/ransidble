package entity

const (
	// ProjectTypeLocal represents a local project
	ProjectTypeLocal = "local"
)

// Project represents a project
type Project struct {
	// Name represents the project name
	Name string `json:"name" validate:"required"`
	// Source represents the project source
	Reference string `json:"reference" validate:"required"`
	// Type represents the project type
	Type string `json:"type" validate:"required"`
}

// GetType returns the project type
func NewProject(name, referene, projectType string) *Project {
	return &Project{
		Name:      name,
		Reference: referene,
		Type:      projectType,
	}
}
