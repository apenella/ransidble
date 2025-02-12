package project

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/model/response"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandle_GetProjectListHandler(t *testing.T) {

	tests := []struct {
		desc               string
		handler            *GetProjectListHandler
		arrangeContextFunc func(r *http.Request, w http.ResponseWriter) echo.Context
		arrangeTestFunc    func(h *GetProjectListHandler)
		assertTestFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			desc: "Testing GetProjectsListHandler.Handle responding with an error when service not initialized and is returning an StatusInternalServerError",
			handler: NewGetProjectListHandler(
				nil,
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body *response.ProjectErrorResponse
				expectedBody := &response.ProjectErrorResponse{
					Error: ErrGetProjectServiceNotInitialized,
				}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "Testing GetProjectsListHandler.Handle project list is returned and is returning an StatusOK",
			handler: NewGetProjectListHandler(
				service.NewMockGetProjectService(),
				logger.NewFakeLogger(),
			),
			arrangeContextFunc: func(r *http.Request, w http.ResponseWriter) echo.Context {
				return echo.New().NewContext(r, w)
			},
			arrangeTestFunc: func(h *GetProjectListHandler) {
				h.service.(*service.MockGetProjectService).On("GetProjectsList").Return([]*entity.Project{
					{
						Name:      "project1",
						Format:    "plain",
						Reference: "project1",
						Storage:   "local",
					},
					{
						Name:      "project2",
						Format:    "plain",
						Reference: "project2",
						Storage:   "local",
					},
				}, nil)
			},
			assertTestFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body []*response.ProjectResponse
				expectedBody := []*response.ProjectResponse{
					{
						Name:      "project1",
						Format:    "plain",
						Reference: "project1",
						Storage:   "local",
					},
					{
						Name:      "project2",
						Format:    "plain",
						Reference: "project2",
						Storage:   "local",
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
			req := httptest.NewRequest(http.MethodGet, "/projects", nil)

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
