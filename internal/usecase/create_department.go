package usecase

import (
	"context"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"unicode/utf8"
)

func (useCase *UseCase) CreateDepartment(ctx context.Context, in in.CreateDepartment) (out.Department, error) {
	if in.Name == "" {
		return out.Department{}, &out.CustomError{Code: 400, Message: "name is required"}
	}

	if utf8.RuneCountInString(in.Name) > 200 {
		return out.Department{}, &out.CustomError{Code: 400, Message: "name length must be between 1 and 200"}
	}

	dep, err := useCase.storage.CreateDepartment(ctx, in.Name, in.ParentId)

	if err != nil {
		return out.Department{}, err
	}

	return out.Department{
		Id:        dep.Id,
		Name:      dep.Name,
		ParentId:  in.ParentId,
		CreatedAt: dep.CreatedAt.Format("02-01-2006"),
	}, nil
}
