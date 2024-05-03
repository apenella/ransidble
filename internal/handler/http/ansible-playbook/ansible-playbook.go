package ansibleplaybook

import (
	"net/http"

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

	ctx := c.Request().Context()

	// Call the business logic to run the Ansible playbook
	err := h.service.Run(ctx, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error running Ansible playbook")
	}

	return c.String(http.StatusOK, "Ansible playbook executed successfully")
}
