package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/model/request"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/test/openapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle_CreateTaskAnsiblePlaybookHandler(t *testing.T) {

	openAPIValidator, err := openapi.PrepareOpenAPIValidator("../../../../api/openapi.yaml")
	if err != nil {
		t.Errorf("Error initializing OpenAPI validator: %s", err)
		t.FailNow()
		return
	}

	tests := []struct {
		desc               string
		handler            *CreateTaskAnsiblePlaybookHandler
		method             string
		path               string
		arrangeContextFunc func(r *http.Request, w http.ResponseWriter) echo.Context
		arrangeTestFunc    func(h *CreateTaskAnsiblePlaybookHandler)
		assertTestFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when project id not provided and is returning a StatusBadRequest",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: ErrProjectIDNotProvided,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when parameters binding fails and is returning a StatusInternalServerError",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {

				requestParameters := &request.AnsiblePlaybookParameters{
					Playbooks: []string{"playbook.yml"},
					Inventory: "inventory.yml",
				}

				body, _ := json.Marshal(requestParameters)
				// The overrided request provides a proper JSON payload but the MIME type is not provided so the binding will fail
				r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))

				c := echo.New().NewContext(r, w)
				c.SetParamNames("project_id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *CreateTaskAnsiblePlaybookHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: fmt.Sprintf("%s: %s", ErrBindingRequestPayload, "code=415, message=Unsupported Media Type"),
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when request payload validation fails and is returning a StatusBadRequest",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("project_id")
				c.SetParamValues("1")
				return c
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					// This is a weak test because it depend on the error message returned by the validation
					Error: fmt.Sprintf("%s: %s", ErrInvalidRequestPayload, "Key: 'AnsiblePlaybookParameters.Playbooks' Error:Field validation for 'Playbooks' failed on the 'required' tag\nKey: 'AnsiblePlaybookParameters.Inventory' Error:Field validation for 'Inventory' failed on the 'required' tag"),
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when receiving and error from the GenerateID method and is returning a StatusInternalServerError",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				requestParameters := &request.AnsiblePlaybookParameters{
					Playbooks: []string{"playbook.yml"},
					Inventory: "inventory.yml",
				}

				body, _ := json.Marshal(requestParameters)
				// The overrided request provides a proper JSON payload. The MIME type is also provided
				r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				c := echo.New().NewContext(r, w)
				c.SetParamNames("project_id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *CreateTaskAnsiblePlaybookHandler) {
				h.service.(*service.MockAnsiblePlaybookService).On("GenerateID").Return("")
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: ErrInvalidTaskID,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when receiving a ProjectNotFoundError error from the Run method and is returning a StatusNotFound",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {

				requestParameters := &request.AnsiblePlaybookParameters{
					Playbooks: []string{"playbook.yml"},
					Inventory: "inventory.yml",
				}

				body, _ := json.Marshal(requestParameters)
				// The overrided request provides a proper JSON payload. The MIME type is also provided
				r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := echo.New().NewContext(r, w)
				c.SetParamNames("project_id")
				c.SetParamValues("1")

				return c
			},
			arrangeTestFunc: func(h *CreateTaskAnsiblePlaybookHandler) {
				h.service.(*service.MockAnsiblePlaybookService).On("GenerateID").Return("testing_task_id")
				h.service.(*service.MockAnsiblePlaybookService).On(
					"Run",
					mock.Anything,
					mock.Anything,
				).Return(
					error.NewProjectNotFoundError(errors.New("testing project not found")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: fmt.Sprintf("%s: %s", ErrRunningAnsiblePlaybook, "testing project not found"),
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		// Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when receiving a ProjectNotProvidedError error from the Run method and is returning a StatusBadRequest
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle responding with an error when receiving a ProjectNotProvidedError error from the Run method and is returning a StatusBadRequest",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {

				requestParameters := &request.AnsiblePlaybookParameters{
					Playbooks: []string{"playbook.yml"},
					Inventory: "inventory.yml",
				}

				body, _ := json.Marshal(requestParameters)
				// The overrided request provides a proper JSON payload. The MIME type is also provided
				r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := echo.New().NewContext(r, w)
				c.SetParamNames("project_id")
				c.SetParamValues("1")

				return c
			},
			arrangeTestFunc: func(h *CreateTaskAnsiblePlaybookHandler) {
				h.service.(*service.MockAnsiblePlaybookService).On("GenerateID").Return("testing_task_id")
				h.service.(*service.MockAnsiblePlaybookService).On("Run", mock.Anything, mock.Anything).Return(
					error.NewProjectNotProvidedError(errors.New("testing project not provided")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: fmt.Sprintf("%s: %s", ErrRunningAnsiblePlaybook, "testing project not provided"),
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing CreateTaskAnsiblePlaybookHandler.Handle succeeded request and is returning a StatusAccepted",
			handler: NewCreateTaskAnsiblePlaybookHandler(
				service.NewMockAnsiblePlaybookService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/tasks/ansible-playbook/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {

				requestParameters := &request.AnsiblePlaybookParameters{
					Playbooks: []string{"playbook.yml"},
					Inventory: "inventory.yml",
				}

				body, _ := json.Marshal(requestParameters)
				// The overrided request provides a proper JSON payload. The MIME type is also provided
				r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := echo.New().NewContext(r, w)
				c.SetParamNames("project_id")
				c.SetParamValues("1")

				return c
			},
			arrangeTestFunc: func(h *CreateTaskAnsiblePlaybookHandler) {
				h.service.(*service.MockAnsiblePlaybookService).On("GenerateID").Return("testing_task_id")
				h.service.(*service.MockAnsiblePlaybookService).On("Run", mock.Anything, mock.Anything).Return(nil)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskCreatedResponse
				expectedBody := &response.TaskCreatedResponse{
					ID: "testing_task_id",
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusAccepted, rec.Code)
			},
		},
	}

	for _, test := range tests {

		rec := httptest.NewRecorder()
		// This is a default request. Depending on the test case the request will be overrided with more specific values
		req := httptest.NewRequest(test.method, test.path, nil)
		context := test.arrangeContextFunc(req, rec)

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrangeTestFunc != nil {
				test.arrangeTestFunc(test.handler)
			}

			err := test.handler.Handle(context)
			assert.NoError(t, err)
			test.assertTestFunc(t, rec)
		})

		t.Run(fmt.Sprintf("OpenAPI %s", test.desc), func(t *testing.T) {
			err := openAPIValidator.ValidateResponse(rec.Body.Bytes(), req, rec.Code, rec.Header())
			assert.NoError(t, err)
		})
	}
}
