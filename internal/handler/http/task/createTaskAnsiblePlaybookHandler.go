package task

import (
	"fmt"
	"net/http"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	request "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

const (
	// ErrInvalidRequestPayload represents an error when the request payload is invalid
	ErrInvalidRequestPayload = "invalid request payload"
	// ErrBindingRequestPayload represents an error when the request payload can not be binded
	ErrBindingRequestPayload = "error binding request payload"
	// ErrStoringTask represents an error when storing a task
	ErrStoringTask = "error storing task"
	// ErrRunningAnsiblePlaybook represents an error when running an ansible playbook
	ErrRunningAnsiblePlaybook = "error running ansible playbook"
	// ErrInvalidTaskID represents an error when the task id is invalid
	ErrInvalidTaskID = "invalid task id"
	// ErrProjectIDNotProvided represents an error when the project id is not provided
	ErrProjectIDNotProvided = "project id not provided"
)

type CreateTaskAnsiblePlaybookHandler struct {
	service service.AnsiblePlaybookServicer
	logger  repository.Logger
}

func NewCreateTaskAnsiblePlaybookHandler(service service.AnsiblePlaybookServicer, logger repository.Logger) *CreateTaskAnsiblePlaybookHandler {
	return &CreateTaskAnsiblePlaybookHandler{
		logger:  logger,
		service: service,
	}
}

func (h *CreateTaskAnsiblePlaybookHandler) Handle(c echo.Context) error {
	var err error
	var errorMsg string
	var res *response.TaskResponse
	var parameters request.AnsiblePlaybookParameters

	ctx := c.Request().Context()

	projectID := c.Param("project_id")
	if projectID == "" {
		res = &response.TaskResponse{
			Error: ErrProjectIDNotProvided,
		}
		h.logger.Error(ErrProjectIDNotProvided)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = c.Bind(&parameters)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrBindingRequestPayload, err.Error())
		res = &response.TaskResponse{
			Error: errorMsg,
		}
		h.logger.Error(errorMsg)
		return c.JSON(http.StatusInternalServerError, res)
	}

	err = parameters.Validate()
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrInvalidRequestPayload, err.Error())
		res = &response.TaskResponse{
			Error: errorMsg,
		}
		h.logger.Error(errorMsg)
		return c.JSON(http.StatusBadRequest, res)
	}

	taskID := h.service.GenerateID()
	if taskID == "" {
		res = &response.TaskResponse{
			Error: ErrInvalidTaskID,
		}
		h.logger.Error(ErrInvalidTaskID)
		return c.JSON(http.StatusInternalServerError, res)
	}

	task := entity.NewTask(taskID, entity.ANSIBLE_PLAYBOOK, &parameters)

	h.logger.Debug(fmt.Sprintf("Creating task %s to run an Ansible playook on project %s\n", taskID, projectID), map[string]interface{}{"component": "handler"})

	err = h.service.Run(ctx, projectID, task)
	if err != nil {
		res = &response.TaskResponse{
			Error: fmt.Sprintf("%s: %s", ErrRunningAnsiblePlaybook, err.Error()),
		}
		return c.JSON(http.StatusServiceUnavailable, res)
	}

	res = &response.TaskResponse{
		ID: taskID,
	}

	return c.JSON(http.StatusAccepted, res)
}
