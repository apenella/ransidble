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
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandle_GetTaskHandler(t *testing.T) {

	tests := []struct {
		desc               string
		handler            *GetTaskHandler
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
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: ErrGetTaskServiceNotInitialized,
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
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *GetTaskHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.TaskErrorResponse
				expectedBody := &response.TaskErrorResponse{
					Error: ErrTaskIDNotProvided,
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
					Error: fmt.Errorf("%s: %s", ErrGettingTask, "testing task not found error").Error(),
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
					Error: fmt.Errorf("%s: %s", ErrGettingTask, "testing task not provided error").Error(),
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
					Error: fmt.Errorf("%s: %s", ErrGettingTask, "testing task unknown error").Error(),
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
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			rec := httptest.NewRecorder()
			// request parameters do not matter for this test. Handler gathers the task id from the context, for this reason, we can use a hardcoded request for all tests
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			context := test.arrangeContextFunc(req, rec)

			if test.arrangeTestFunc != nil {
				test.arrangeTestFunc(test.handler)
			}

			err := test.handler.Handle(context)
			assert.NoError(t, err)

			test.assertTestFunc(t, rec)
		})
	}
}
