package project

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

func TestHandle_GetProjectHandler(t *testing.T) {

	openAPIValidator, err := openapi.PrepareOpenAPIValidator("../../../../api/openapi.yaml")
	if err != nil {
		t.Errorf("Error initializing OpenAPI validator: %s", err)
		t.FailNow()
		return
	}

	tests := []struct {
		desc               string
		handler            *GetProjectHandler
		arrangeContextFunc func(r *http.Request, w http.ResponseWriter) echo.Context
		arrangeTestFunc    func(h *GetProjectHandler)
		assertTestFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			desc: "Testing GetProjectHandler.Handle responding with an error when service not initialized and is returning an StatusInternalServerError",
			handler: NewGetProjectHandler(
				nil,
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  ErrGetProjectServiceNotInitialized,
					Status: http.StatusInternalServerError,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing GetProjectHandler.Handle responding with an error when project id not provided and is returning an StatusBadRequest",
			handler: NewGetProjectHandler(
				service.NewMockGetProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *GetProjectHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  ErrProjectIDNotProvided,
					Status: http.StatusBadRequest,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing GetProjectHandler.Handle responding with an error when project not found and is returning an StatusNotFound",
			handler: NewGetProjectHandler(
				service.NewMockGetProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetProjectHandler) {
				h.service.(*service.MockGetProjectService).On("GetProject", "1").Return(
					&entity.Project{},
					error.NewProjectNotFoundError(errors.New("testing project not found error")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Errorf("%s: %s", ErrGettingProject, "testing project not found error").Error(),
					Status: http.StatusNotFound,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			desc: "Testing GetProjectHandler.Handle responding with an error when reciving a project not provided error and is returning an StatusBadRequest",
			handler: NewGetProjectHandler(
				service.NewMockGetProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetProjectHandler) {
				h.service.(*service.MockGetProjectService).On("GetProject", "1").Return(
					&entity.Project{},
					error.NewProjectNotProvidedError(errors.New("testing project not provided error")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Errorf("%s: %s", ErrGettingProject, "testing project not provided error").Error(),
					Status: http.StatusBadRequest,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing GetProjectHandler.Handle responding with an error when gets a project unknown error and is returning an StatusInternalServerError",
			handler: NewGetProjectHandler(
				service.NewMockGetProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetProjectHandler) {
				h.service.(*service.MockGetProjectService).On("GetProject", "1").Return(
					&entity.Project{},
					errors.New("testing project unknown error"),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Errorf("%s: %s", ErrGettingProject, "testing project unknown error").Error(),
					Status: http.StatusInternalServerError,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing GetProjectHandler.Handle request success and is returning an StatusOK",
			handler: NewGetProjectHandler(
				service.NewMockGetProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("1")
				return c
			},
			arrangeTestFunc: func(h *GetProjectHandler) {
				h.service.(*service.MockGetProjectService).On("GetProject", "1").Return(
					&entity.Project{
						Format:    "plain",
						Name:      "project1",
						Reference: "project1",
						Storage:   "local",
					},
					nil,
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectResponse
				expectedBody := &response.ProjectResponse{
					Format:    "plain",
					Name:      "project1",
					Reference: "project1",
					Storage:   "local",
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
	}

	for _, test := range tests {
		var req *http.Request
		rec := httptest.NewRecorder()

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			// request parameters do not matter for this test. Handler gathers the project id from the context, for this reason, we can use a hardcoded request for all tests
			req = httptest.NewRequest(http.MethodGet, "/projects/project-1", nil)
			context := test.arrangeContextFunc(req, rec)

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
