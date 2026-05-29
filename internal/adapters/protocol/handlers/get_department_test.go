package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUseCaseService struct {
	mock.Mock
}

func (m *MockUseCaseService) CreateDepartment(ctx context.Context, in in.CreateDepartment) (out.Department, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(out.Department), args.Error(1)
}

func (m *MockUseCaseService) CreateEmployee(ctx context.Context, parentId string, in in.CreateEmployee) (out.Employee, error) {
	args := m.Called(ctx, parentId, in)
	return args.Get(0).(out.Employee), args.Error(1)
}

func (m *MockUseCaseService) DeleteDepartment(ctx context.Context, id string, mode string, reassignToDepartmentId *string) error {
	args := m.Called(ctx, id, mode, reassignToDepartmentId)
	return args.Error(0)
}

func (m *MockUseCaseService) GetDepartment(ctx context.Context, id string, depth string, includeEmployees string) (out.DepartmentSingleInfo, error) {
	args := m.Called(ctx, id, depth, includeEmployees)

	if args.Get(0) == nil {
		return out.DepartmentSingleInfo{}, args.Error(1)
	}

	return *args.Get(0).(*out.DepartmentSingleInfo), args.Error(1)
}

func (m *MockUseCaseService) UpdateDepartment(ctx context.Context, id string, in in.UpdateDepartment) (out.DepartmentSingleInfo, error) {
	args := m.Called(ctx, id, in)
	return args.Get(0).(out.DepartmentSingleInfo), args.Error(1)
}

func TestHGetDepartment200(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	mockTaskService.On("GetDepartment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&out.DepartmentSingleInfo{}, nil)

	req := httptest.NewRequest("GET", "/departments/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HGetDepartment(mockTaskService)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	slog.Info("Success", slog.Any("success", w.Body.String()))

	mockTaskService.AssertExpectations(t)
}

func TestHGetDepartment404(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	mockTaskService.On("GetDepartment", mock.Anything, "1", mock.Anything, mock.Anything).Return(nil,
		&out.CustomError{
			Code:    404,
			Message: "Department is not found",
		})

	req := httptest.NewRequest("GET", "/departments/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HGetDepartment(mockTaskService)(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	slog.Info("Error", slog.Any("error", w.Body.String()))
	mockTaskService.AssertExpectations(t)
}

func TestHGetDepartment500(t *testing.T) {
	mockTaskService := new(MockUseCaseService)

	mockTaskService.On("GetDepartment", mock.Anything, "2", mock.Anything, mock.Anything).Return(nil,
		errors.New("unexpected error"))

	req := httptest.NewRequest("GET", "/departments/2", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	HGetDepartment(mockTaskService)(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	slog.Info("Error", slog.Any("error", w.Body.String()))
	mockTaskService.AssertExpectations(t)
}
