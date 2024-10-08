package project

import (
	"errors"
	"fmt"
	"net/http"

	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/mapper"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

// GetProjectHandler struct to handle requests to get a project details
type GetProjectHandler struct {
	service service.GetProjectServicer
	logger  repository.Logger
}

// NewGetProjectHandler creates a new GetProjectHandler
func NewGetProjectHandler(s service.GetProjectServicer, logger repository.Logger) *GetProjectHandler {
	return &GetProjectHandler{
		service: s,
		logger:  logger,
	}
}

// Handle method to get a project details
func (h *GetProjectHandler) Handle(c echo.Context) error {

	var errorMsg string
	var errorResponse *response.ProjectErrorResponse
	var httpStatus int
	var projectNotFoundErr *domainerror.ProjectNotFoundError
	var projectNotProvidedErr *domainerror.ProjectNotProvidedError

	if h.service == nil {
		errorResponse = &response.ProjectErrorResponse{
			Error: ErrGetProjectServiceNotInitialized,
		}

		h.logger.Error(ErrGetProjectServiceNotInitialized, map[string]interface{}{
			"component": "GetProjectHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
		})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	id := c.Param("id")
	if id == "" {
		errorResponse = &response.ProjectErrorResponse{
			Error: ErrProjectIDNotProvided,
		}
		h.logger.Error(ErrProjectIDNotProvided, map[string]interface{}{
			"component": "GetProjectHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
		})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	h.logger.Debug(
		fmt.Sprintf("getting project %s\n", id),
		map[string]interface{}{
			"component":  "GetProjectHandler.Handle",
			"package":    "github.com/apenella/ransidble/internal/handler/http/project",
			"project_id": id,
		})

	project, err := h.service.GetProject(id)
	if err != nil {
		httpStatus = http.StatusInternalServerError

		if errors.As(err, &projectNotFoundErr) {
			httpStatus = http.StatusNotFound
		}

		if errors.As(err, &projectNotProvidedErr) {
			httpStatus = http.StatusBadRequest
		}

		errorMsg = fmt.Sprintf("%s: %s", ErrGettingProject, err.Error())

		errorResponse = &response.ProjectErrorResponse{
			Error: errorMsg,
		}

		h.logger.Error(errorMsg, map[string]interface{}{
			"component":  "GetProjectHandler.Handle",
			"package":    "github.com/apenella/ransidble/internal/handler/http/project",
			"project_id": id,
		})
		return c.JSON(httpStatus, errorResponse)
	}

	projectMapper := mapper.NewProjectMapper()
	projectResponse := projectMapper.ToProjectResponse(project)

	return c.JSON(http.StatusOK, projectResponse)
}
