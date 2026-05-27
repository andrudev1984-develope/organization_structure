package usecase

import (
	"context"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"strconv"
	"time"
	"unicode/utf8"
)

func (useCase *UseCase) CreateEmployee(ctx context.Context, parentId string, in in.CreateEmployee) (out.Employee, error) {
	if in.FullName == "" {
		return out.Employee{}, &out.CustomError{Code: 400, Message: "full name is required"}
	}

	if utf8.RuneCountInString(in.FullName) > 200 {
		return out.Employee{}, &out.CustomError{Code: 400, Message: "full name length must be between 1 and 200"}
	}

	if in.Position == "" {
		return out.Employee{}, &out.CustomError{Code: 400, Message: "position is required"}
	}

	if utf8.RuneCountInString(in.Position) > 200 {
		return out.Employee{}, &out.CustomError{Code: 400, Message: "position length must be between 1 and 200"}
	}

	hiredAt, err := mustParse(in.HiredAt)

	if err != nil {
		return out.Employee{}, &out.CustomError{Code: 400, Message: "hiredAt date must be in 'DD-MM-YYYY' format'"}
	}

	depId, err := strconv.Atoi(parentId)

	if err != nil {
		return out.Employee{}, &out.CustomError{Code: 400, Message: "Invalid department id format"}
	}

	emp, err := useCase.storage.CreateEmployee(ctx, in.FullName, uint(depId), in.Position, hiredAt)

	if err != nil {
		return out.Employee{}, err
	}

	return out.Employee{
		Id:        emp.Id,
		FullName:  emp.FullName,
		Position:  emp.Position,
		HiredAt:   formatDate((*time.Time)(emp.HiredAt)),
		CreatedAt: emp.CreatedAt.Format("02-01-2006"),
	}, nil
}

func formatDate(date *time.Time) *string {
	if date == nil {
		return nil
	}

	return new((*date).Format("02-01-2006"))
}

func mustParse(date *string) (*time.Time, error) {
	if date == nil {
		return nil, nil
	}

	t, e := time.Parse("02-01-2006", *date)

	if e != nil {
		return nil, e
	}

	return new(t), nil
}
