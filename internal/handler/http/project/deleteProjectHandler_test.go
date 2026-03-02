package project

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/test/openapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandle_DeleteProjectHandler(t *testing.T) {
	openAPIValidator, err := openapi.PrepareOpenAPIValidator("../../../../api/openapi.yaml")
	if err != nil {
		t.Errorf("Error initializing OpenAPI validator: %s", err)
		t.FailNow()
		return
	}

	tests := []struct {
		desc               string
		handler            *DeleteProjectHandler
		path               string
		arrangeContextFunc func(r *http.Request, w http.ResponseWriter) echo.Context
		arrangeTestFunc    func(t *testing.T, h *DeleteProjectHandler)
		assertTestFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			desc: "Testing DeleteProjectHandler.Handle responding with an error when service not initialized and is returning an StatusInternalServerError",
			handler: NewDeleteProjectHandler(
				nil,
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			path: "/projects/test-id",
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  ErrDeleteProjectServiceNotInitialized,
					Status: http.StatusInternalServerError,
				}

				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing DeleteProjectHandler.Handle responding with an error when project id not provided and is returning an StatusBadRequest",
			handler: NewDeleteProjectHandler(
				service.NewMockDeleteProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			path:            "/projects/test-id", // test-id not used in the test, you need to set it in the context to make the test work
			arrangeTestFunc: func(t *testing.T, h *DeleteProjectHandler) {},
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
			desc: "Testing DeleteProjectHandler.Handle responding with an error when project not found and is returning an StatusInternalServerError",
			handler: NewDeleteProjectHandler(
				service.NewMockDeleteProjectService(),
				logger.NewFakeLogger(),
			),
			path: "/projects/test-id", // test-id not used in the test, you need to set it in the context to make the test work
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("test-id")
				return c
			},
			arrangeTestFunc: func(t *testing.T, h *DeleteProjectHandler) {
				h.service.(*service.MockDeleteProjectService).On(
					"Delete",
					"test-id",
				).Return(
					fmt.Errorf("project not found"),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrDeletingProject, "project not found"),
					Status: http.StatusInternalServerError,
				}

				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing DeleteProjectHandler.Handle responding with no content when project is not found and is returning an StatusNotFound",
			handler: NewDeleteProjectHandler(
				service.NewMockDeleteProjectService(),
				logger.NewFakeLogger(),
			),
			path: "/projects/test-id", // test-id not used in the test, you need to set it in the context to make the test work
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("test-id")
				return c
			},
			arrangeTestFunc: func(t *testing.T, h *DeleteProjectHandler) {
				h.service.(*service.MockDeleteProjectService).On(
					"Delete",
					"test-id",
				).Return(
					domainerror.NewProjectNotFoundError(fmt.Errorf("project not found")),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			desc: "Testing DeleteProjectHandler.Handle responding with no content when project is deleted successfully and is returning an StatusNoContent",
			handler: NewDeleteProjectHandler(
				service.NewMockDeleteProjectService(),
				logger.NewFakeLogger(),
			),
			path: "/projects/test-id", // test-id not used in the test, you need to set it in the context to make the test work
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				c := echo.New().NewContext(r, w)
				c.SetParamNames("id")
				c.SetParamValues("test-id")
				return c
			},
			arrangeTestFunc: func(t *testing.T, h *DeleteProjectHandler) {
				h.service.(*service.MockDeleteProjectService).On(
					"Delete",
					"test-id",
				).Return(
					nil,
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		// This is a default request. Depending on the test case the request will be overrided with more specific values
		req := httptest.NewRequest(http.MethodDelete, test.path, nil)
		context := test.arrangeContextFunc(req, rec)

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrangeTestFunc != nil {
				test.arrangeTestFunc(t, test.handler)
			}

			err := test.handler.Handle(context)
			assert.NoError(t, err)

			test.assertTestFunc(t, rec)
		})

		_ = openAPIValidator
		t.Run(fmt.Sprintf("OpenAPI %s", test.desc), func(t *testing.T) {
			err := openAPIValidator.ValidateResponse(rec.Body.Bytes(), req, rec.Code, rec.Header())
			assert.NoError(t, err)
		})
	}

}
