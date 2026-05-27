package usecase

import (
	"context"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"strconv"
	"unicode/utf8"
)

func (useCase *UseCase) UpdateDepartment(ctx context.Context, id string, in in.UpdateDepartment) (out.DepartmentSingleInfo, error) {
	pId, err := strconv.Atoi(id)

	if err != nil {
		return out.DepartmentSingleInfo{}, &out.CustomError{
			Code:    400,
			Message: "Need valid department id",
		}
	}

	if in.Name != nil {
		if *in.Name == "" {
			return out.DepartmentSingleInfo{}, &out.CustomError{Code: 400, Message: "Department name is required"}
		}

		if utf8.RuneCountInString(*in.Name) > 200 {
			return out.DepartmentSingleInfo{}, &out.CustomError{Code: 400, Message: "Department name length must be between 1 and 200"}
		}
	}

	if in.ParentId != nil && *in.ParentId == uint(pId) {
		return out.DepartmentSingleInfo{}, &out.CustomError{Code: 400,
			Message: "Department cannot be child of itself"}
	}

	err = useCase.storage.UpdateDepartment(ctx, uint(pId), in.Name, in.ParentId)

	if err != nil {
		return out.DepartmentSingleInfo{}, err
	}

	return useCase.GetDepartment(ctx, id, strconv.Itoa(1), "true")
}
