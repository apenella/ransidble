package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/test/openapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandle_GetTaskHandler(t *testing.T) {

	openAPIValidator, err := openapi.PrepareOpenAPIValidator("../../../../api/openapi.yaml")
	if err != nil {
		t.Errorf("Error initializing OpenAPI validator: %s", err)
		t.FailNow()
		return
	}

	tests := []struct {
		desc               string
		handler            *GetTaskHandler
		method             string
		path               string
		arrangeContextFunc func(r *http.Request, w http.ResponseWriter) echo.Context
		arrangeTestFunc    func(h *GetTaskHandler)
		assertTestFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			desc: "Testing GetTaskHandler.Handle responding with an error when service not initialized and is returning an StatusInternalServerError",
			handler: NewGetTaskHandler(
				nil,
				logger.NewFakeLogger(),
			),
			method: http.MethodGet,
			path:   "/tasks/task-id",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error:  ErrGetTaskServiceNotInitialized,
					Status: http.StatusInternalServerError,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing GetTaskHandler.Handle responding with an error when task id not provided and is returning an StatusBadRequest",
			handler: NewGetTaskHandler(
				service.NewMockGetTaskService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodGet,
			path:   "/tasks/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *GetTaskHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error:  ErrTaskIDNotProvided,
					Status: http.StatusBadRequest,
				}

				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing GetTaskHandler.Handle responding with an error when task not found and is returning an StatusNotFound",
			handler: NewGetTaskHandler(
				service.NewMockGetTaskService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodGet,
			path:   "/tasks/task-id",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetTaskHandler) {
				h.service.(*service.MockGetTaskService).On("GetTask", "1").Return(
					&entity.Task{},
					error.NewTaskNotFoundError(errors.New("testing task not found error")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error:  fmt.Errorf("%s: %s", ErrGettingTask, "testing task not found error").Error(),
					Status: http.StatusNotFound,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			desc: "Testing GetTaskHandler.Handle responding with an error when receiving a task not provided error and is returning an StatusBadRequest",
			handler: NewGetTaskHandler(
				service.NewMockGetTaskService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodGet,
			path:   "/tasks/1",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetTaskHandler) {
				h.service.(*service.MockGetTaskService).On("GetTask", "1").Return(
					&entity.Task{},
					error.NewTaskNotProvidedError(errors.New("testing task not provided error")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error:  fmt.Errorf("%s: %s", ErrGettingTask, "testing task not provided error").Error(),
					Status: http.StatusBadRequest,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing GetTaskHandler.Handle responding with an error when gets a task unknown error and is returning an StatusInternalServerError",
			handler: NewGetTaskHandler(
				service.NewMockGetTaskService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodGet,
			path:   "/tasks/task-id",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetTaskHandler) {
				h.service.(*service.MockGetTaskService).On("GetTask", "1").Return(
					&entity.Task{},
					errors.New("testing task unknown error"),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error:  fmt.Errorf("%s: %s", ErrGettingTask, "testing task unknown error").Error(),
					Status: http.StatusInternalServerError,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing GetTaskHandler.Handle request success and is returning an StatusOK",
			handler: NewGetTaskHandler(
				service.NewMockGetTaskService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodGet,
			path:   "/tasks/task-id",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetTaskHandler) {
				h.service.(*service.MockGetTaskService).On("GetTask", "1").Return(
					&entity.Task{
						ID:        "1",
						ProjectID: "project1",
						Command:   entity.AnsiblePlaybookCommand,
						Parameters: &entity.AnsiblePlaybookParameters{
							Playbooks: []string{"playbook.yml"},
							Inventory: "inventory.yml",
						},
						CompletedAt: "0000-01-01T01:01:01",
						CreatedAt:   "0000-01-01T01:01:01",
						ExecutedAt:  "0000-01-01T01:01:01",
						Status:      entity.ACCEPTED,
					},
					nil,
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskResponse
				expectedBody := &response.TaskResponse{
					ID:        "1",
					ProjectID: "project1",
					Command:   entity.AnsiblePlaybookCommand,
					Parameters: map[string]interface{}{
						"playbooks": []interface{}{"playbook.yml"},
						"inventory": "inventory.yml",
					},
					CompletedAt: "0000-01-01T01:01:01",
					CreatedAt:   "0000-01-01T01:01:01",
					ExecutedAt:  "0000-01-01T01:01:01",
					Status:      entity.ACCEPTED,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		rec := httptest.NewRecorder()

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
