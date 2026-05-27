package usecase

import (
	"context"
	"organization_structure/internal/model"
	"time"
)

type IStorage interface {
	CreateDepartment(ctx context.Context, name string, parentId *uint) (*model.Department, error)
	GetDepartment(ctx context.Context, departmentId uint, depth int, includeEmployees bool) (*model.Department, error)
	UpdateDepartment(ctx context.Context, departmentId uint, name *string, parentId *uint) error
	DeleteDepartment(ctx context.Context, departmentId uint, mode string, reassignToDepartmentId *uint) error
	CreateEmployee(ctx context.Context, fullName string, departmentId uint, position string, hiredAt *time.Time) (*model.Employee, error)
}

type UseCase struct {
	storage IStorage
}

func NewUseCase(storage IStorage) *UseCase {
	return &UseCase{storage: storage}
}
