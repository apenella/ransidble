package entity

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	// ProjectTypeLocal represents a local project
	ProjectTypeLocal = "local"
	// ProjectFormatPlain represents project in plain format
	ProjectFormatPlain = "plain"
	// ProjectFormatTarGz represents a project in tar.gz format
	ProjectFormatTarGz = "targz"

	// ExtensionTarGz represents the tar.gz extension
	ExtensionTarGz = ".tar.gz"
)

// Project entity represents a project
type Project struct {
	// Format represents the project format. This field is required and must be one of the following values: plain, targz
	Format string `json:"format" validate:"required,oneof=plain targz"`
	// Name represents the project name. This field is required
	Name string `json:"name" validate:"required"`
	// Source represents the project source. This field is required
	Reference string `json:"reference" validate:"required"`
	// Storage represents the project type. This field is required and must be one of the following values: local
	Storage string `json:"storage" validate:"required,oneof=local"`
	// Hash represents the project hash
	//Hash string `json:"hash" validate:"sha256"`
}

// NewProject creates a new project instance
func NewProject(name, referene, format, storage string) *Project {
	return &Project{
		Format:    format,
		Name:      name,
		Reference: referene,
		Storage:   storage,
	}
}

// Validate validates the project entity
func (p *Project) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ValidateProjectFormat validates the project format
func ValidateProjectFormat(format string) error {
	validate := validator.New()
	err := validate.Var(format, "required,oneof=plain targz")
	if err != nil {
		return fmt.Errorf("Invalid format: %s. %w", format, err)
	}

	return nil
}

// ValidateProjectStorage validates the project storage
func ValidateProjectStorage(storage string) error {
	validate := validator.New()
	err := validate.Var(storage, "required,oneof=local")

	if err != nil {
		return fmt.Errorf("Invalid storage type: %s. %w", storage, err)
	}

	return nil
}

// ValidateProjectFileExtension validates the project file extension
func ValidateProjectFileExtension(file string) error {
	has := strings.HasSuffix(file, ExtensionTarGz)
	if !has {
		return fmt.Errorf("file %s extension not supported", file)
	}
	return nil
}
