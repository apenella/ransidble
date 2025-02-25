package openapi

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

const (
	// OpenAPIDefPath is the path to the OpenAPI definition
	OpenAPIDefPath = "../../api/openapi.yaml"
)

// Validator is a struct to validate the OpenAPI specification
type Validator struct {
	loader *openapi3.Loader
	router routers.Router
}

// PrepareOpenAPIValidator creates a new OpenAPIValidator
func PrepareOpenAPIValidator(openAPIDefPath string) (*Validator, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(openAPIDefPath)
	if err != nil {
		return nil, fmt.Errorf("error loading OpenAPI definition: %w", err)
	}

	err = doc.Validate(loader.Context)
	if err != nil {
		return nil, fmt.Errorf("error validating OpenAPI definition: %w", err)
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("error creating router: %w", err)
	}

	return &Validator{
		loader: loader,
		router: router,
	}, nil
}

// ValidateResponse validates the response against the OpenAPI specification
func (v *Validator) ValidateResponse(responseBody []byte, request *http.Request, statusCode int, headers http.Header) error {

	route, pathParams, err := v.router.FindRoute(request)
	if err != nil {
		return fmt.Errorf("error finding route: %w", err)
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    request,
		PathParams: pathParams,
		Route:      route,
	}

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 statusCode,
		Header:                 headers,
		Options: &openapi3filter.Options{
			IncludeResponseStatus: true,
		},
	}

	responseValidationInput.SetBodyBytes(responseBody)
	err = openapi3filter.ValidateResponse(v.loader.Context, responseValidationInput)
	if err != nil {
		return fmt.Errorf("response validation failed: %w", err)
	}

	return nil
}
