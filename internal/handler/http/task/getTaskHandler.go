package task

import (
	"fmt"
	"net/http"

	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

const (
	// ErrGetTaskServiceNotInitialized represents an error when the GetTaskService is not initialized
	ErrGetTaskServiceNotInitialized = "get task service not initialized"
	// ErrTaskIDNotProvided represents an error when the task id is not provided
	ErrTaskIDNotProvided = "task id not provided"
)

type GetTaskHandler struct {
	service service.GetTaskServicer
	logger  repository.Logger
}

func NewGetTaskHandler(s service.GetTaskServicer, logger repository.Logger) *GetTaskHandler {
	return &GetTaskHandler{
		service: s,
		logger:  logger,
	}
}

func (h *GetTaskHandler) Handle(c echo.Context) error {

	var res *response.CommandResponse

	if h.service == nil {
		res = &response.CommandResponse{
			Error: ErrGetTaskServiceNotInitialized,
		}

		h.logger.Error(ErrGetTaskServiceNotInitialized)
		return c.JSON(http.StatusInternalServerError, res)
	}

	id := c.Param("id")
	if id == "" {
		res = &response.CommandResponse{
			Error: ErrTaskIDNotProvided,
		}
		h.logger.Error(ErrTaskIDNotProvided)
		return c.JSON(http.StatusBadRequest, res)
	}

	h.logger.Debug(fmt.Sprintf("getting task %s\n", id), map[string]interface{}{"component": "handler"})
	task, err := h.service.GetTask(id)
	if err != nil {
		res = &response.CommandResponse{
			Error: fmt.Sprintf("%s: %s", ErrGetTaskServiceNotInitialized, err.Error()),
		}
		h.logger.Error(fmt.Sprintf("%s: %s", ErrGetTaskServiceNotInitialized, err.Error()))
		return c.JSON(http.StatusNotFound, res)
	}

	return c.JSON(http.StatusNotImplemented, task)
}
