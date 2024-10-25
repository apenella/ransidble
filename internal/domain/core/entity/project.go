package entity

import "github.com/go-playground/validator/v10"

const (
	// ProjectTypeLocal represents a local project
	ProjectTypeLocal = "local"
	// ProjectFormatPlain represents project in plain format
	ProjectFormatPlain = "plain"
	// ProjectFormatTarGz represents a project in tar.gz format
	ProjectFormatTarGz = "targz"
)

// Project entity represents a project
type Project struct {
	// Format represents the project format. This field is required and must be one of the following values: plain, targz
	Format string `json:"format" validate:"required,oneof=plain targz"`
	// Name represents the project name. This field is required
	Name string `json:"name" validate:"required"`
	// Source represents the project source. This field is required
	Reference string `json:"reference" validate:"required"`
	// Type represents the project type. This field is required and must be one of the following values: local
	Type string `json:"type" validate:"required,oneof=local"`
}

// NewProject creates a new project instance
func NewProject(name, referene, format, projectType string) *Project {
	return &Project{
		Format:    format,
		Name:      name,
		Reference: referene,
		Type:      projectType,
	}
}

// Validate validates the project entity
func (p *Project) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
