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

func TestHCreateEmployee201(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	mockTaskService.On("CreateEmployee", mock.Anything, mock.Anything, mock.Anything).Return(out.Employee{}, nil)

	req := httptest.NewRequest("POST", "/departments/1/employees", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HCreateEmployee(mockTaskService)(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	slog.Info("Success", slog.Any("Success", w.Body.String()))

	mockTaskService.AssertExpectations(t)
}
