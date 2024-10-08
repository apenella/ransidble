package entity

const (
	// ProjectTypeLocal represents a local project
	ProjectTypeLocal = "local"
	// ProjectFormatPlain represents project in plain format
	ProjectFormatPlain = "plain"
	// ProjectFormatTarGz represents a project in tar.gz format
	ProjectFormatTarGz = "targz"
)

// Project represents a project
type Project struct {
	// Format represents the project format
	Format string `json:"format" validate:"required"`
	// Name represents the project name
	Name string `json:"name" validate:"required"`
	// Source represents the project source
	Reference string `json:"reference" validate:"required"`
	// Type represents the project type
	Type string `json:"type" validate:"required"`
}

// NewProject creates a new project
func NewProject(name, referene, format, projectType string) *Project {
	return &Project{
		Format:    format,
		Name:      name,
		Reference: referene,
		Type:      projectType,
	}
}
