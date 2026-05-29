package handlers

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"organization_structure/internal/model/out"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHUpdateDepartment400(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	req := httptest.NewRequest("PATH", "/departments/1", bytes.NewBufferString("Invalid json body"))
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HUpdateDepartment(mockTaskService)(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	slog.Info("Error", slog.Any("Error", w.Body.String()))

	mockTaskService.AssertExpectations(t)
}

func TestHUpdateDepartment200(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	mockTaskService.On("UpdateDepartment", mock.Anything, mock.Anything, mock.Anything).Return(out.DepartmentSingleInfo{}, nil)

	req := httptest.NewRequest("PATCH", "/departments/1", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HUpdateDepartment(mockTaskService)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	slog.Info("Success", slog.Any("Success", w.Body.String()))

	mockTaskService.AssertExpectations(t)
}
