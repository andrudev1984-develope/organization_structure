package usecase

import (
	"log/slog"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateDepartmentUseCase(test *testing.T) {
	var useCase = NewUseCase(MockDb{})

	var cases = []struct {
		caseName string
		id       string
		in       in.UpdateDepartment
		isError  bool
		out      out.Department
	}{
		{
			caseName: "Id format error",
			id:       "",
			in:       in.UpdateDepartment{},
			isError:  true,
		},
		{
			caseName: "Name need error",
			id:       "1",
			in:       in.UpdateDepartment{Name: new(string)},
			isError:  true,
		},
		{
			caseName: "Name length error",
			id:       "1",
			in:       in.UpdateDepartment{Name: new(strings.Repeat("a", 201))},
			isError:  true,
		},
		{
			caseName: "Department cannot be child of itself error",
			id:       "1",
			in:       in.UpdateDepartment{Name: new("Name"), ParentId: new(uint(1))},
			isError:  true,
		},
		{
			caseName: "Success",
			id:       "1",
			in:       in.UpdateDepartment{Name: new("Name"), ParentId: new(uint(2))},
			isError:  false,
		},
	}

	for _, testCase := range cases {
		test.Run(testCase.caseName, func(test *testing.T) {
			var res, err = useCase.UpdateDepartment(test.Context(), testCase.id, testCase.in)
			if testCase.isError {
				assert.Error(test, err)
				slog.Error("Error operation", slog.Any("err", err.Error()))
			} else {
				assert.NoError(test, err)
				slog.Info("Success operation", slog.Any("id", res.Department.Id), slog.Any("name", res.Department.Name),
					slog.Any("created_at", res.Department.CreatedAt))
			}
		})
	}
}
