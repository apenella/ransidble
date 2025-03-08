package request

import "github.com/go-playground/validator/v10"

// ProjectParameters represents a request describing a project
type ProjectParameters struct {
	// Format represents the project format
	Format string `json:"format" validate:"required,oneof=targz plain"`
	// Name represents the project name
	// Name string `json:"name" validate:"required"`
	// // Source represents the project source
	// Reference string `json:"reference" validate:"required"`
	// Storage represents the project type
	Storage string `json:"storage" validate:"required,oneof=local"`
}

// Validate validates the request
func (p *ProjectParameters) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
