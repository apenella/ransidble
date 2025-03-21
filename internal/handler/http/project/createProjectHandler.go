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
	var requestParameters request.ProjectParameters
	var errorResponse *response.ProjectErrorResponse
	var projectFileHeader *multipart.FileHeader
	var projectReceivedFile multipart.File

	var projectAlreadyExists *domainerror.ProjectAlreadyExistsError

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

	projectFileHeader, err = c.FormFile(RequestFormProjectFileFieldeName)
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

	projectReceivedFile, err = projectFileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("error opening file: %s", err.Error()))
	}

	err = h.service.Create(requestParameters.Format, requestParameters.Storage, projectFileHeader.Filename, projectReceivedFile)
	if err != nil {

		httpStatus := http.StatusInternalServerError

		if errors.As(err, &projectAlreadyExists) {
			httpStatus = http.StatusConflict
		}

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
		return c.JSON(httpStatus, errorResponse)
	}

	// projectSourceFile, err = projectFileHeader.Open()
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
	// 			"file":      projectFileHeader.Filename,
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
	// 			"file":        projectFileHeader.Filename,
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
	// 			"file":        projectFileHeader.Filename,
	// 			"destination": projectDestinationFile.Name(),
	// 		})
	// 	return c.JSON(http.StatusInternalServerError, errorResponse)
	// }

	c.JSON(http.StatusOK, requestParameters)

	return nil
}
