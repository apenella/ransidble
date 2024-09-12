package task

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/mapper"
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
	var httpStatus int
	var requestParameters request.AnsiblePlaybookParameters
	var projectNotFoundErr *domainerror.ProjectNotFoundError
	var projectNotProvidedErr *domainerror.ProjectNotProvidedError
	var errorResponse *response.TaskErrorResponse

	ctx := c.Request().Context()

	projectID := c.Param("project_id")
	if projectID == "" {
		errorResponse = &response.TaskErrorResponse{
			Error: ErrProjectIDNotProvided,
		}
		h.logger.Error(
			ErrProjectIDNotProvided,
			map[string]interface{}{
				"component": "CreateTaskAnsiblePlaybookHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/task",
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	err = c.Bind(&requestParameters)
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrBindingRequestPayload, err.Error())
		errorResponse = &response.TaskErrorResponse{
			Error: errorMsg,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":  "CreateTaskAnsiblePlaybookHandler.Handle",
				"package":    "github.com/apenella/ransidble/internal/handler/http/task",
				"project_id": projectID,
			})
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	err = requestParameters.Validate()
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", ErrInvalidRequestPayload, err.Error())
		errorResponse = &response.TaskErrorResponse{
			Error: errorMsg,
		}
		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":  "CreateTaskAnsiblePlaybookHandler.Handle",
				"package":    "github.com/apenella/ransidble/internal/handler/http/task",
				"project_id": projectID,
			})
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	ansiblePlaybookParametersMapper := mapper.NewAnsiblePlaybookParametersMapper()
	parameters := ansiblePlaybookParametersMapper.ToAnsiblePlaybookParametersEntity(&requestParameters)

	taskID := h.service.GenerateID()
	if taskID == "" {
		errorResponse = &response.TaskErrorResponse{
			Error: ErrInvalidTaskID,
		}
		h.logger.Error(
			ErrInvalidTaskID,
			map[string]interface{}{
				"component": "CreateTaskAnsiblePlaybookHandler.Handle",
				"package":   "github.com/apenella/ransidble/internal/handler/http/task",
				"task_id":   taskID,
			})

		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	task := entity.NewTask(taskID, entity.ANSIBLE_PLAYBOOK, parameters)

	h.logger.Debug(
		fmt.Sprintf("creating task %s to run an Ansible playbook on project %s\n", taskID, projectID),
		map[string]interface{}{
			"component":  "CreateTaskAnsiblePlaybookHandler.Handle",
			"package":    "github.com/apenella/ransidble/internal/handler/http/task",
			"project_id": projectID,
			"task_id":    taskID,
		})

	err = h.service.Run(ctx, projectID, task)
	if err != nil {
		httpStatus = http.StatusInternalServerError

		if errors.As(err, &projectNotFoundErr) {
			httpStatus = http.StatusNotFound
		}

		if errors.As(err, &projectNotProvidedErr) {
			httpStatus = http.StatusBadRequest
		}

		errorMsg = fmt.Sprintf("%s: %s", ErrRunningAnsiblePlaybook, err.Error())
		errorResponse = &response.TaskErrorResponse{
			Error: errorMsg,
		}

		h.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":  "CreateTaskAnsiblePlaybookHandler.Handle",
				"package":    "github.com/apenella/ransidble/internal/handler/http/task",
				"project_id": projectID,
				"task_id":    taskID,
			})

		return c.JSON(httpStatus, errorResponse)
	}

	taskCreated := &response.TaskCreatedResponse{
		ID: taskID,
	}

	return c.JSON(http.StatusAccepted, taskCreated)
}
