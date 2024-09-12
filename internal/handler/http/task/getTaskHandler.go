package task

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
	// ErrGetTaskServiceNotInitialized represents an error when the GetTaskService is not initialized
	ErrGetTaskServiceNotInitialized = "get task service not initialized"
	// ErrTaskIDNotProvided represents an error when the task id is not provided
	ErrTaskIDNotProvided = "task id not provided"
	// ErrGettingTask represents an error executing the method getting task
	ErrGettingTask = "error getting task"
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

	var errorResponse *response.TaskErrorResponse
	var errorMsg string
	var httpStatus int
	var taskNotFoundErr *domainerror.TaskNotFoundError
	var taskNotProvidedErr *domainerror.TaskNotProvidedError

	if h.service == nil {
		errorResponse = &response.TaskErrorResponse{
			Error: ErrGetTaskServiceNotInitialized,
		}

		h.logger.Error(
			ErrGetTaskServiceNotInitialized,
			map[string]interface{}{
				"component": "GetTaskHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/task",
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	id := c.Param("id")
	if id == "" {
		errorResponse = &response.TaskErrorResponse{
			Error: ErrTaskIDNotProvided,
		}
		h.logger.Error(
			ErrTaskIDNotProvided,
			map[string]interface{}{
				"component": "GetTaskHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/task",
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	h.logger.Debug(
		fmt.Sprintf("getting task %s\n", id),
		map[string]interface{}{
			"component": "GetTaskHandler.Handle",
			"package":   "github.com/apenella/ransidble/internal/handler/http/task",
			"task_id":   id,
		})
	task, err := h.service.GetTask(id)
	if err != nil {

		httpStatus = http.StatusInternalServerError

		if errors.As(err, &taskNotFoundErr) {
			httpStatus = http.StatusNotFound
		}

		if errors.As(err, &taskNotProvidedErr) {
			httpStatus = http.StatusBadRequest
		}

		errorMsg = fmt.Sprintf("%s: %s", ErrGettingTask, err.Error())

		errorResponse = &response.TaskErrorResponse{
			Error: errorMsg,
		}

		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component": "GetTaskHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/task",
				"task_id":   id,
			})
		return c.JSON(httpStatus, errorResponse)
	}

	taskMapper := mapper.NewTaskMapper()
	taskResponse := taskMapper.ToTaskResponse(task)

	return c.JSON(http.StatusOK, taskResponse)
}
