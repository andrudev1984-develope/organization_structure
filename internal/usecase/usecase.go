package usecase

import (
	"context"
	"organization_structure/internal/model"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"time"
)

type IStorage interface {
	CreateDepartment(ctx context.Context, name string, parentId *uint) (*model.Department, error)
	GetDepartment(ctx context.Context, departmentId uint, depth int, includeEmployees bool) (*model.Department, error)
	UpdateDepartment(ctx context.Context, departmentId uint, name *string, parentId *uint) error
	DeleteDepartment(ctx context.Context, departmentId uint, mode string, reassignToDepartmentId *uint) error
	CreateEmployee(ctx context.Context, fullName string, departmentId uint, position string, hiredAt *time.Time) (*model.Employee, error)
}

type IUSeCase interface {
	CreateDepartment(ctx context.Context, in in.CreateDepartment) (out.Department, error)
	CreateEmployee(ctx context.Context, parentId string, in in.CreateEmployee) (out.Employee, error)
	DeleteDepartment(ctx context.Context, id string, mode string, reassignToDepartmentId *string) error
	GetDepartment(ctx context.Context, id string, depth string, includeEmployees string) (out.DepartmentSingleInfo, error)
	UpdateDepartment(ctx context.Context, id string, in in.UpdateDepartment) (out.DepartmentSingleInfo, error)
}

type UseCase struct {
	storage IStorage
}

func NewUseCase(storage IStorage) *UseCase {
	return &UseCase{storage: storage}
}
