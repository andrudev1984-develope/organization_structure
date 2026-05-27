package usecase

import (
	"log/slog"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEmployeeUseCase(test *testing.T) {
	var useCase = NewUseCase(MockDb{})

	var cases = []struct {
		caseName string
		parentId string
		in       in.CreateEmployee
		isError  bool
		out      out.Department
	}{
		{caseName: "Need full name error",
			parentId: "1",
			in: in.CreateEmployee{
				FullName: "",
				Position: "",
				HiredAt:  nil,
			}, isError: true, out: out.Department{}},
		{caseName: "Full name length error",
			parentId: "1",
			in: in.CreateEmployee{
				FullName: strings.Repeat("a", 201),
				Position: "",
				HiredAt:  nil,
			}, isError: true, out: out.Department{}},
		{caseName: "Need position error",
			parentId: "1",
			in: in.CreateEmployee{
				FullName: "1",
				Position: "",
				HiredAt:  nil,
			}, isError: true, out: out.Department{}},
		{caseName: "Position length error",
			parentId: "1",
			in: in.CreateEmployee{
				FullName: "1",
				Position: strings.Repeat("a", 201),
				HiredAt:  nil,
			}, isError: true, out: out.Department{}},
		{caseName: "Hired at error",
			parentId: "1",
			in: in.CreateEmployee{
				FullName: "1",
				Position: "1",
				HiredAt:  new("invalid"),
			}, isError: true, out: out.Department{}},
		{caseName: "Parent id error",
			parentId: "1x",
			in: in.CreateEmployee{
				FullName: "1",
				Position: "1",
				HiredAt:  new("31-12-2025"),
			}, isError: true, out: out.Department{}},
		{caseName: "Parent is not found",
			parentId: "0",
			in: in.CreateEmployee{
				FullName: "1",
				Position: "1",
				HiredAt:  new("31-12-2025"),
			}, isError: true, out: out.Department{}},
		{caseName: "Success",
			parentId: "1",
			in: in.CreateEmployee{
				FullName: "1",
				Position: "1",
				HiredAt:  new("31-12-2025"),
			}, isError: false, out: out.Department{}},
	}

	for _, testCase := range cases {
		test.Run(testCase.caseName, func(test *testing.T) {
			var res, err = useCase.CreateEmployee(test.Context(), testCase.parentId, testCase.in)
			if testCase.isError {
				assert.Error(test, err)
				slog.Error("Error operation", slog.Any("err", err.Error()))
			} else {
				assert.NoError(test, err)
				slog.Info("Success operation", slog.Any("id", res.Id), slog.Any("name", res.FullName),
					slog.Any("created_at", res.CreatedAt))
			}
		})
	}
}
