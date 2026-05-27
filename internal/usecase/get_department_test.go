package usecase

import (
	"log/slog"
	"organization_structure/internal/model/out"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDepartmentUseCase(test *testing.T) {
	var useCase = NewUseCase(MockDb{})

	var cases = []struct {
		caseName         string
		id               string
		depth            string
		includeEmployees string
		isError          bool
		out              out.Department
	}{
		{
			caseName:         "Id format error",
			id:               "x",
			depth:            "x",
			includeEmployees: "x",
			isError:          true,
		},
		{
			caseName:         "Depth format error",
			id:               "1",
			depth:            "x",
			includeEmployees: "x",
			isError:          true,
		},
		{
			caseName:         "Depth constraint error",
			id:               "1",
			depth:            "10",
			includeEmployees: "x",
			isError:          true,
		},
		{
			caseName:         "IncludeEmployees format error",
			id:               "1",
			depth:            "5",
			includeEmployees: "x",
			isError:          true,
		},
		{
			caseName:         "Success",
			id:               "1",
			depth:            "5",
			includeEmployees: "true",
			isError:          false,
		},
	}

	for _, testCase := range cases {
		test.Run(testCase.caseName, func(test *testing.T) {
			var res, err = useCase.GetDepartment(test.Context(), testCase.id, testCase.depth, testCase.includeEmployees)
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
