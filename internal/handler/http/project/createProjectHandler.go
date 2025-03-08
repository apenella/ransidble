package project

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/apenella/ransidble/internal/domain/core/model/request"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

const (
	// RequestFormProjectMetadataFieldName represents the form field name for the project metadata
	RequestFormProjectMetadataFieldName = "metadata"
	// RequestFormProjectFileFieldeName represents the form field name for the project file
	RequestFormProjectFileFieldeName = "file"

	// ErrReadingFormProjectMetadataField represents an error when the form field for the project metadata can not be read
	ErrReadingFormProjectMetadataField = "error reading project metadata field"
	// ErrReadingFormProjectFileField represents an error when the form field for the project file can not be read
	ErrReadingFormProjectFileField = "error reading project file field"
	// ErrInvalidRequestMetadata represents an error when the request metadata is invalid
	ErrInvalidRequestMetadata = "invalid request metadata"
	// ErrCreatingProject represents an error when the project can not be created
	ErrCreatingProject = "error creating project"
)

// CreateProjectHandler handles the request to create a new project
type CreateProjectHandler struct {
	service service.CreateProjectServicer
	logger  repository.Logger
}

// NewCreateProjectHandler creates a new CreateProjectHandler
func NewCreateProjectHandler(service service.CreateProjectServicer, logger repository.Logger) *CreateProjectHandler {
	return &CreateProjectHandler{
		service: service,
		logger:  logger,
	}
}

// Handle method to create a new project
func (h *CreateProjectHandler) Handle(c echo.Context) error {
	var err error
	var errorMsg string
	var requestParameters request.ProjectParameters
	var errorResponse *response.ProjectErrorResponse
	var projectFieldHeader *multipart.FileHeader

	metadata := c.FormValue(RequestFormProjectMetadataFieldName)
	err = json.Unmarshal([]byte(metadata), &requestParameters)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrReadingFormProjectMetadataField, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error: errorMsg,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	err = requestParameters.Validate()
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrInvalidRequestMetadata, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error: errorMsg,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	projectFieldHeader, err = c.FormFile(RequestFormProjectFileFieldeName)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrReadingFormProjectFileField, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error: errorMsg,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	err = h.service.Create(requestParameters.Format, requestParameters.Storage, projectFieldHeader)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrCreatingProject, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error: errorMsg,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	// projectSourceFile, err = projectFieldHeader.Open()
	// if err != nil {
	// 	errorMsg = fmt.Sprintf("%s: %s", "error opening source file", err.Error())
	// 	errorResponse = &response.ProjectErrorResponse{
	// 		Error: errorMsg,
	// 	}
	// 	h.logger.Error(
	// 		errorMsg,
	// 		map[string]interface{}{
	// 			"component": "CreateProjectHandler.Handle",
	// 			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
	// 			"file":      projectFieldHeader.Filename,
	// 		})
	// 	return c.JSON(http.StatusInternalServerError, errorResponse)
	// }
	// defer projectSourceFile.Close()

	// projectDestinationFile, err = os.OpenFile("tmpFile.tar.gz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	errorMsg = fmt.Sprintf("%s: %s", "error opening destination file", err.Error())
	// 	errorResponse = &response.ProjectErrorResponse{
	// 		Error: errorMsg,
	// 	}
	// 	h.logger.Error(
	// 		errorMsg,
	// 		map[string]interface{}{
	// 			"component":   "CreateProjectHandler.Handle",
	// 			"package":     "github.com/apenella/ransidble/internal/handler/http/project",
	// 			"file":        projectFieldHeader.Filename,
	// 			"destination": projectDestinationFile.Name(),
	// 		})
	// 	return c.JSON(http.StatusInternalServerError, errorResponse)
	// }
	// defer projectDestinationFile.Close()

	// _, err = io.Copy(projectDestinationFile, projectSourceFile)
	// if err != nil {
	// 	errorMsg = fmt.Sprintf("%s: %s", "error copying source file to destination file", err.Error())
	// 	errorResponse = &response.ProjectErrorResponse{
	// 		Error: errorMsg,
	// 	}
	// 	h.logger.Error(
	// 		errorMsg,
	// 		map[string]interface{}{
	// 			"component":   "CreateProjectHandler.Handle",
	// 			"package":     "github.com/apenella/ransidble/internal/handler/http/project",
	// 			"file":        projectFieldHeader.Filename,
	// 			"destination": projectDestinationFile.Name(),
	// 		})
	// 	return c.JSON(http.StatusInternalServerError, errorResponse)
	// }

	c.JSON(http.StatusOK, requestParameters)

	return nil
}
