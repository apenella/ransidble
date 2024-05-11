package ansibleplaybook

import (
	"net/http"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	model "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/labstack/echo/v4"
)

type AnsiblePlaybookHandler struct {
	service service.AnsiblePlaybookServicer
}

func NewAnsiblePlaybookHandler(service service.AnsiblePlaybookServicer) *AnsiblePlaybookHandler {
	return &AnsiblePlaybookHandler{
		service: service,
	}
}

func (h *AnsiblePlaybookHandler) Handle(c echo.Context) error {
	var err error

	ctx := c.Request().Context()

	// Call the business logic to run the Ansible playbook
	// Aquí cal inidcar que es crida a l'executor de l'ansible-playbook
	// - aquest pot ser un executor genèric que depenent la Task instància un executor concret

	// The handler should get from the post request a data structure with the playbook details to run and set it in to AnsiblePlaybookOptions

	var options model.AnsiblePlaybookParameters

	// Here we should get the data from the post request and set it in the options variable
	err = c.Bind(&options)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	// I need to create a new Task instance with the options as parameters
	id := h.service.GenerateID()
	task := entity.NewTask(id, entity.ANSIBLE_PLAYBOOK, options)

	err = h.service.Run(ctx, task)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error running Ansible playbook")
	}

	return c.JSON(http.StatusAccepted, task)
}
