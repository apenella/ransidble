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
	Source string `json:"source" validate:"required"`
	// Type represents the project type
	Type string `json:"type" validate:"required"`
}

// GetType returns the project type
func NewProject(name, source, projectType string) *Project {
	return &Project{
		Name:   name,
		Source: source,
		Type:   projectType,
	}
}
