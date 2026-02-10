package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/model/request"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	serverhttp "github.com/apenella/ransidble/internal/handler/http"
	"github.com/labstack/echo/v4"
)

const (
	// RequestFormProjectMetadataFieldName represents the form field name for the project metadata
	RequestFormProjectMetadataFieldName = "metadata"
	// RequestFormProjectFileFieldeName represents the form field name for the project file
	RequestFormProjectFileFieldeName = "file"
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
	var errorResponse *response.ProjectErrorResponse
	var projectAlreadyExists *domainerror.ProjectAlreadyExistsError
	var projectErrorResponseStatus int
	var projectFileHeader *multipart.FileHeader
	var projectID string
	var projectReceivedFile multipart.File
	var requestParameters request.ProjectParameters

	if h.service == nil {
		errorResponse = &response.ProjectErrorResponse{
			Error:  ErrCreateProjectServiceNotInitialized,
			Status: http.StatusInternalServerError,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	metadata := c.FormValue(RequestFormProjectMetadataFieldName)
	err = json.Unmarshal([]byte(metadata), &requestParameters)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrReadingFormProjectMetadataField, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error:  errorMsg,
			Status: http.StatusInternalServerError,
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
			Error:  errorMsg,
			Status: http.StatusBadRequest,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	projectFileHeader, err = c.FormFile(RequestFormProjectFileFieldeName)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrReadingFormProjectFileField, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error:  errorMsg,
			Status: http.StatusBadRequest,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	projectReceivedFile, err = projectFileHeader.Open()
	if err != nil {
		errorMsg = fmt.Sprintf("error opening file: %s", err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error:  errorMsg,
			Status: http.StatusInternalServerError,
		}

		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})

		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	projectID, err = h.service.Create(requestParameters.Format, requestParameters.Storage, projectFileHeader.Filename, projectReceivedFile)
	if err != nil {

		httpStatus := http.StatusInternalServerError
		projectErrorResponseStatus = http.StatusInternalServerError

		if errors.As(err, &projectAlreadyExists) {
			httpStatus = http.StatusConflict
			projectErrorResponseStatus = http.StatusConflict
		}

		errorMsg = fmt.Sprintf("%s: %s", ErrCreatingProject, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error:  errorMsg,
			Status: projectErrorResponseStatus,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "CreateProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(httpStatus, errorResponse)
	}

	// Use the route constant but replace the parameter placeholder with the actual ID
	location := fmt.Sprintf("%s/%s", serverhttp.CreateProjectPath, projectID)
	c.Response().Header().Set("Location", location)

	return c.NoContent(http.StatusCreated)
}
