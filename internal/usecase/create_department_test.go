package usecase

import (
	"context"
	"log/slog"
	"organization_structure/internal/model"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

type MockDb struct {
}

func (m MockDb) CreateDepartment(ctx context.Context, name string, parentId *uint) (*model.Department, error) {
	if parentId != nil && *parentId == 0 {
		return nil, out.CustomError{
			Code:    404,
			Message: "Parent department is not found",
		}
	}

	return &model.Department{
		Id:        0,
		ParentID:  new(uint),
		Name:      "name",
		CreatedAt: time.Now(),
	}, nil
}

func (m MockDb) GetDepartment(ctx context.Context, departmentId uint, depth int, includeEmployees bool) (*model.Department, error) {
	return &model.Department{
		Id:        departmentId,
		Name:      "Name",
		CreatedAt: time.Now(),
	}, nil
}

func (m MockDb) UpdateDepartment(ctx context.Context, departmentId uint, name *string, parentId *uint) error {
	return nil
}

func (m MockDb) DeleteDepartment(ctx context.Context, departmentId uint, mode string, reassignToDepartmentId *uint) error {
	return nil
}

func (m MockDb) CreateEmployee(ctx context.Context, fullName string, departmentId uint, position string, hiredAt *time.Time) (*model.Employee, error) {
	if departmentId == 0 {
		return nil, out.CustomError{
			Code:    404,
			Message: "Parent department is not found",
		}
	}

	return &model.Employee{
		Id:           0,
		DepartmentID: departmentId,
		FullName:     fullName,
		Position:     position,
		HiredAt:      (*datatypes.Date)(hiredAt),
		CreatedAt:    time.Now(),
	}, nil
}

func TestCreateUseCase(test *testing.T) {
	var useCase = NewUseCase(MockDb{})

	var cases = []struct {
		caseName string
		in       in.CreateDepartment
		isError  bool
		out      out.Department
	}{
		{caseName: "Need name error",
			in: in.CreateDepartment{
				Name:     "",
				ParentId: nil,
			}, isError: true, out: out.Department{}},
		{caseName: "Name length error",
			in: in.CreateDepartment{
				Name:     strings.Repeat("a", 201),
				ParentId: nil,
			}, isError: true, out: out.Department{}},
		{caseName: "Parent is not found",
			in: in.CreateDepartment{
				Name:     "Name",
				ParentId: new(uint),
			}, isError: true, out: out.Department{}},
		{caseName: "Success",
			in: in.CreateDepartment{
				Name:     "Name",
				ParentId: nil,
			}, isError: false, out: out.Department{}},
	}

	for _, testCase := range cases {
		test.Run(testCase.caseName, func(test *testing.T) {
			var res, err = useCase.CreateDepartment(test.Context(), testCase.in)
			if testCase.isError {
				assert.Error(test, err)
				slog.Error("Error operation", slog.Any("err", err.Error()))
			} else {
				assert.NoError(test, err)
				slog.Info("Success operation", slog.Any("id", res.Id), slog.Any("name", res.Name),
					slog.Any("created_at", res.CreatedAt))
			}
		})
	}
}
