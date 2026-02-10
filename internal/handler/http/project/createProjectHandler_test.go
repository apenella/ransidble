package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/core/model/request"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/test/openapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle_CreateProjectHandler(t *testing.T) {

	openAPIValidator, err := openapi.PrepareOpenAPIValidator("../../../../api/openapi.yaml")
	if err != nil {
		t.Errorf("Error initializing OpenAPI validator: %s", err)
		t.FailNow()
		return
	}

	tests := []struct {
		desc               string
		handler            *CreateProjectHandler
		method             string
		path               string
		arrangeContextFunc func(r *http.Request, w http.ResponseWriter) echo.Context
		arrangeTestFunc    func(h *CreateProjectHandler)
		assertTestFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when service is not initialized and is returning a StatusInternalServerError",
			handler: NewCreateProjectHandler(
				nil,
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  ErrCreateProjectServiceNotInitialized,
					Status: http.StatusInternalServerError,
				}

				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when unmarshaling project metadata form fails and is returning a StatusInternalServerError",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				// The error in this test case if forced by sending an empty body

				r = httptest.NewRequest(http.MethodPost, "/projects", nil)
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {

				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrReadingFormProjectMetadataField, "unexpected end of JSON input"),
					Status: http.StatusInternalServerError,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when project metadata is set with an invalid format and is returning a StatusBadRequest",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				// The erorr in test case is forced by setting an invalid format in the request

				var bodyBuffer bytes.Buffer

				requestParameters := &request.ProjectParameters{
					Format:  "invalid-format",
					Storage: entity.ProjectTypeLocal,
				}

				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					t.Fatal(err)
				}

				multiparWriter := multipart.NewWriter(&bodyBuffer)
				defer multiparWriter.Close()

				multiparWriter.WriteField(RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				r = httptest.NewRequest(http.MethodPost, "/projects", &bodyBuffer)
				r.Header.Set(echo.HeaderContentType, multiparWriter.FormDataContentType())
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrInvalidRequestMetadata, "Key: 'ProjectParameters.Format' Error:Field validation for 'Format' failed on the 'oneof' tag"),
					Status: http.StatusBadRequest,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when project metadata is set with an invalid storage and is returning a StatusBadRequest",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				// The erorr in test case is forced by setting an invalid storage in the request

				var bodyBuffer bytes.Buffer

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: "invalid-storage",
				}

				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					t.Fatal(err)
				}

				multiparWriter := multipart.NewWriter(&bodyBuffer)
				defer multiparWriter.Close()

				multiparWriter.WriteField(RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				r = httptest.NewRequest(http.MethodPost, "/projects", &bodyBuffer)
				r.Header.Set(echo.HeaderContentType, multiparWriter.FormDataContentType())
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrInvalidRequestMetadata, "Key: 'ProjectParameters.Storage' Error:Field validation for 'Storage' failed on the 'oneof' tag"),
					Status: http.StatusBadRequest,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when is was not possible achieve the project file form the form and is returning a StatusBadRequest",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				var bodyBuffer bytes.Buffer

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: entity.ProjectTypeLocal,
				}

				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					t.Fatal(err)
				}

				multiparWriter := multipart.NewWriter(&bodyBuffer)
				defer multiparWriter.Close()

				multiparWriter.WriteField(RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				r = httptest.NewRequest(http.MethodPost, "/projects", &bodyBuffer)
				r.Header.Set(echo.HeaderContentType, multiparWriter.FormDataContentType())
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrReadingFormProjectFileField, "http: no such file"),
					Status: http.StatusBadRequest,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when is was not possible to create the project and is returning a StatusInternalServerError",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				var bodyBuffer bytes.Buffer

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: entity.ProjectTypeLocal,
				}

				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					t.Fatal(err)
				}

				multiparWriter := multipart.NewWriter(&bodyBuffer)
				defer multiparWriter.Close()

				multiparWriter.WriteField(RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				part, err := multiparWriter.CreateFormFile(RequestFormProjectFileFieldeName, "project.tar.gz")
				if err != nil {
					t.Fatal(err)
				}
				projectContentFile := strings.NewReader("project-content")
				_, err = io.Copy(part, projectContentFile)
				if err != nil {
					t.Fatal(err)
				}

				r = httptest.NewRequest(http.MethodPost, "/projects", &bodyBuffer)
				r.Header.Set(echo.HeaderContentType, multiparWriter.FormDataContentType())
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {
				h.service.(*service.MockCreateProjectService).On(
					"Create",
					entity.ProjectFormatTarGz,
					entity.ProjectTypeLocal,
					mock.Anything,
				).Return("no-project-id", fmt.Errorf("error opening project file"))
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrCreatingProject, "error opening project file"),
					Status: http.StatusInternalServerError,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle responding with an error when the project already exists and is returning a StatusConflict",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				var bodyBuffer bytes.Buffer

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: entity.ProjectTypeLocal,
				}

				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					t.Fatal(err)
				}

				multiparWriter := multipart.NewWriter(&bodyBuffer)
				defer multiparWriter.Close()

				multiparWriter.WriteField(RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				part, err := multiparWriter.CreateFormFile(RequestFormProjectFileFieldeName, "project.tar.gz")
				if err != nil {
					t.Fatal(err)
				}
				projectContentFile := strings.NewReader("project-content")
				_, err = io.Copy(part, projectContentFile)
				if err != nil {
					t.Fatal(err)
				}

				r = httptest.NewRequest(http.MethodPost, "/projects", &bodyBuffer)
				r.Header.Set(echo.HeaderContentType, multiparWriter.FormDataContentType())
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {
				h.service.(*service.MockCreateProjectService).On(
					"Create",
					entity.ProjectFormatTarGz,
					entity.ProjectTypeLocal,
					mock.Anything,
				).Return(
					"no-project-id",
					domainerror.NewProjectAlreadyExistsError(
						fmt.Errorf("project already exists"),
					),
				)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error:  fmt.Sprintf("%s: %s", ErrCreatingProject, "project already exists"),
					Status: http.StatusConflict,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusConflict, rec.Code)
			},
		},
		{
			desc: "Testing CreateProjectHandler.Handle request success and it is returning a StatusCreated",
			handler: NewCreateProjectHandler(
				service.NewMockCreateProjectService(),
				logger.NewFakeLogger(),
			),
			method: http.MethodPost,
			path:   "/projects",
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				var bodyBuffer bytes.Buffer

				requestParameters := &request.ProjectParameters{
					Format:  entity.ProjectFormatTarGz,
					Storage: entity.ProjectTypeLocal,
				}

				requestParametersJSON, err := json.Marshal(requestParameters)
				if err != nil {
					t.Fatal(err)
				}

				multiparWriter := multipart.NewWriter(&bodyBuffer)
				defer multiparWriter.Close()

				multiparWriter.WriteField(RequestFormProjectMetadataFieldName, string(requestParametersJSON))

				part, err := multiparWriter.CreateFormFile(RequestFormProjectFileFieldeName, "project.tar.gz")
				if err != nil {
					t.Fatal(err)
				}
				projectContentFile := strings.NewReader("project-content")
				_, err = io.Copy(part, projectContentFile)
				if err != nil {
					t.Fatal(err)
				}

				r = httptest.NewRequest(http.MethodPost, "/projects", &bodyBuffer)
				r.Header.Set(echo.HeaderContentType, multiparWriter.FormDataContentType())
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *CreateProjectHandler) {
				h.service.(*service.MockCreateProjectService).On(
					"Create",
					entity.ProjectFormatTarGz,
					entity.ProjectTypeLocal,
					mock.Anything,
				).Return("project-id", nil)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Equal(t, rec.Header().Get("Location"), "/projects/project-id")
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
