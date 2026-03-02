package project

import (
	"fmt"
	"net/http"

	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

// DeleteProjectHandler is the HTTP handler for deleting a project.
type DeleteProjectHandler struct {
	service service.DeleteProjectServicer
	logger  repository.Logger
}

// NewDeleteProjectHandler creates a new instance of DeleteProjectHandler.
func NewDeleteProjectHandler(service service.DeleteProjectServicer, logger repository.Logger) *DeleteProjectHandler {
	return &DeleteProjectHandler{
		service: service,
		logger:  logger,
	}
}

// Handle handles the HTTP request for deleting a project.
func (h *DeleteProjectHandler) Handle(c echo.Context) error {
	var err error
	var errorMsg string
	var errorResponse *response.ProjectErrorResponse
	var projectErrorResponseStatus int
	var projectID string

	if h.service == nil {
		errorResponse = &response.ProjectErrorResponse{
			Error:  ErrDeleteProjectServiceNotInitialized,
			Status: http.StatusInternalServerError,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "DeleteProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	projectID = c.Param("id")

	if len(projectID) == 0 {
		errorResponse = &response.ProjectErrorResponse{
			Error:  ErrProjectIDNotProvided,
			Status: http.StatusBadRequest,
		}
		h.logger.Error(
			ErrProjectIDNotProvided,
			map[string]interface{}{
				"component": "DeleteProjectHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/project",
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	err = h.service.Delete(projectID)
	if err != nil {
		if _, ok := err.(*domainerror.ProjectNotFoundError); ok {
			projectErrorResponseStatus = http.StatusNotFound
		} else {
			projectErrorResponseStatus = http.StatusInternalServerError
		}
		errorMsg = fmt.Sprintf("%s: %s", ErrDeletingProject, err.Error())
		errorResponse = &response.ProjectErrorResponse{
			Error:  errorMsg,
			Status: projectErrorResponseStatus,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":  "DeleteProjectHandler.Handle",
				"package":    "github.com/apenella/ransidble/internal/handler/http/project",
				"project_id": projectID,
			})
		return c.JSON(projectErrorResponseStatus, errorResponse)
	}

	return c.NoContent(http.StatusNoContent)
}
