package usecase

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteDepartmentUseCase(test *testing.T) {
	var useCase = NewUseCase(MockDb{})

	var cases = []struct {
		caseName               string
		id                     string
		mode                   string
		reassignToDepartmentId *string
		isError                bool
	}{
		{
			caseName:               "Id format error",
			id:                     "e",
			mode:                   "",
			reassignToDepartmentId: nil,
			isError:                true,
		},
		{
			caseName:               "Mode error",
			id:                     "1",
			mode:                   "mode",
			reassignToDepartmentId: nil,
			isError:                true,
		},
		{
			caseName:               "Need reassignToDepartmentId",
			id:                     "1",
			mode:                   "reassign",
			reassignToDepartmentId: nil,
			isError:                true,
		},
		{
			caseName:               "ReassignToDepartmentId format error",
			id:                     "1",
			mode:                   "reassign",
			reassignToDepartmentId: new("x"),
			isError:                true,
		},
		{
			caseName:               "Success",
			id:                     "1",
			mode:                   "reassign",
			reassignToDepartmentId: new("1"),
			isError:                false,
		},
	}

	for _, testCase := range cases {
		test.Run(testCase.caseName, func(test *testing.T) {
			var err = useCase.DeleteDepartment(test.Context(), testCase.id, testCase.mode, testCase.reassignToDepartmentId)
			if testCase.isError {
				assert.Error(test, err)
				slog.Error("Error operation", slog.Any("err", err.Error()))
			} else {
				assert.NoError(test, err)
			}
		})
	}
}
