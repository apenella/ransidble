package ansibleplaybook

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
)

type AnsiblePlaybookHandler struct {
	service   service.AnsiblePlaybookServicer
	taskStore repository.TaskStorer
	logger    repository.Logger
}

func NewAnsiblePlaybookHandler(service service.AnsiblePlaybookServicer, taskStore repository.TaskStorer, logger repository.Logger) *AnsiblePlaybookHandler {
	return &AnsiblePlaybookHandler{
		logger:    logger,
		service:   service,
		taskStore: taskStore,
	}
}

func (h *AnsiblePlaybookHandler) Handle(c echo.Context) error {
	var err error
	var errorMsg string
	var res *response.CommandResponse
	var parameters request.AnsiblePlaybookParameters

	ctx := c.Request().Context()

	err = c.Bind(&parameters)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrBindingRequestPayload, err.Error())
		res = &response.CommandResponse{
			Error: errorMsg,
		}
		h.logger.Error(errorMsg)
		return c.JSON(http.StatusInternalServerError, res)
	}

	err = parameters.Validate()
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrInvalidRequestPayload, err.Error())
		res = &response.CommandResponse{
			Error: errorMsg,
		}
		h.logger.Error(errorMsg)
		return c.JSON(http.StatusBadRequest, res)
	}

	id := h.service.GenerateID()
	if id == "" {
		res = &response.CommandResponse{
			Error: ErrInvalidTaskID,
		}
		h.logger.Error(ErrInvalidTaskID)
		return c.JSON(http.StatusInternalServerError, res)
	}

	task := entity.NewTask(id, entity.ANSIBLE_PLAYBOOK, &parameters)
	h.logger.Debug(fmt.Sprintf("running task %s\n", id), map[string]interface{}{"component": "handler", "task": task})

	err = h.service.Run(ctx, task)
	if err != nil {
		res = &response.CommandResponse{
			Error: fmt.Sprintf("%s: %s", ErrRunningAnsiblePlaybook, err.Error()),
		}
		return c.JSON(http.StatusServiceUnavailable, res)
	}

	res = &response.CommandResponse{
		ID: id,
	}

	return c.JSON(http.StatusAccepted, res)
}
