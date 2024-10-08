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

const (
// ErrGetProjectServiceNotInitialized represents an error when the GetProjectService is not initialized
)

// GetProjecListtHandler struct to handle get project requests
type GetProjecListtHandler struct {
	service service.GetProjectServicer
	logger  repository.Logger
}

// NewGetProjecListtHandler creates a new GetProjecListtHandler
func NewGetProjecListtHandler(s service.GetProjectServicer, logger repository.Logger) *GetProjecListtHandler {
	return &GetProjecListtHandler{
		service: s,
		logger:  logger,
	}
}

// Handle method to get a task
func (h *GetProjecListtHandler) Handle(c echo.Context) error {

	var errorMsg string
	var errorResponse *response.ProjectErrorResponse
	var httpStatus int
	var projectNotFoundErr *domainerror.ProjectNotFoundError

	if h.service == nil {
		errorResponse = &response.ProjectErrorResponse{
			Error: ErrGetProjectServiceNotInitialized,
		}

		h.logger.Error(ErrGetProjectServiceNotInitialized, map[string]interface{}{
			"component": "GetProjecListtHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
		})

		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	h.logger.Debug("getting project list", map[string]interface{}{
		"component": "GetProjecListtHandler.Handle",
		"package":   "github.com/apenella/ransidble/internal/handler/http/project",
	})

	projects, err := h.service.GetProjectsList()
	if err != nil {
		httpStatus = http.StatusInternalServerError

		if errors.As(err, &projectNotFoundErr) {
			httpStatus = http.StatusNotFound
		}

		errorMsg = fmt.Sprintf("%s: %s", ErrGettingProjectList, err.Error())

		h.logger.Error(errorMsg, map[string]interface{}{
			"component": "GetProjecListtHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/project",
		})

		errorResponse = &response.ProjectErrorResponse{
			Error: errorMsg,
		}

		return c.JSON(httpStatus, errorResponse)
	}

	projectListResponse := make([]*response.ProjectResponse, 0)
	projectMapper := mapper.NewProjectMapper()
	for _, project := range projects {
		projectListResponse = append(projectListResponse, projectMapper.ToProjectResponse(project))
	}

	return c.JSON(http.StatusOK, projectListResponse)
}
