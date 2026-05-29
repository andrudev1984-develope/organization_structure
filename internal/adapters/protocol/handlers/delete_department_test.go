package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHDeleteDepartment204(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	mockTaskService.On("DeleteDepartment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest("DELETE", "/departments/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HDeleteDepartment(mockTaskService)(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Nil(t, w.Body.Bytes())

	mockTaskService.AssertExpectations(t)
}
