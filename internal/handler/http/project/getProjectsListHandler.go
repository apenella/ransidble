package project

import (
	"fmt"
	"net/http"

	"github.com/apenella/ransidble/internal/domain/core/mapper"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

const (
// ErrGetProjectServiceNotInitialized represents an error when the GetProjectService is not initialized
)

// GetProjectListHandler struct to handle get project requests
type GetProjectListHandler struct {
	service service.GetProjectServicer
	logger  repository.Logger
}

// NewGetProjectListHandler creates a new GetProjectListHandler
func NewGetProjectListHandler(s service.GetProjectServicer, logger repository.Logger) *GetProjectListHandler {
	return &GetProjectListHandler{
		service: s,
		logger:  logger,
	}
}

// Handle method to get a task
func (h *GetProjectListHandler) Handle(c echo.Context) error {

	var errorMsg string
	var errorResponse *response.ProjectErrorResponse

	if h.service == nil {
		errorResponse = &response.ProjectErrorResponse{
			Error:  ErrGetProjectServiceNotInitialized,
			Status: http.StatusInternalServerError,
		}

		h.logger.Error(ErrGetProjectServiceNotInitialized, map[string]interface{}{
			"component": "GetProjectListHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
		})

		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	h.logger.Debug("getting project list", map[string]interface{}{
		"component": "GetProjectListHandler.Handle",
		"package":   "github.com/apenella/ransidble/internal/handler/http/project",
	})

	projects, err := h.service.GetProjectsList()
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrGettingProjectList, err.Error())

		h.logger.Error(errorMsg, map[string]interface{}{
			"component": "GetProjectListHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
		})

		errorResponse = &response.ProjectErrorResponse{
			Error:  errorMsg,
			Status: http.StatusInternalServerError,
		}

		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	projectListResponse := make([]*response.ProjectResponse, 0)
	projectMapper := mapper.NewProjectMapper()
	for _, project := range projects {
		projectListResponse = append(projectListResponse, projectMapper.ToProjectResponse(project))
	}

	return c.JSON(http.StatusOK, projectListResponse)
}
