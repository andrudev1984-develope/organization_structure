package usecase

import (
	"context"
	"organization_structure/internal/model/out"
	"strconv"
)

func (useCase *UseCase) DeleteDepartment(ctx context.Context, id string, mode string, reassignToDepartmentId *string) error {
	pId, err := strconv.Atoi(id)

	if err != nil {
		return &out.CustomError{
			Code:    400,
			Message: "Need valid department id",
		}
	}

	if mode != "cascade" && mode != "reassign" {
		return &out.CustomError{
			Code:    400,
			Message: "Need delete mode : cascade or reassign",
		}
	}

	if mode == "reassign" && (reassignToDepartmentId == nil || *reassignToDepartmentId == "") {
		return &out.CustomError{
			Code:    400,
			Message: "Need reassign_to_department_id for reassign mode",
		}
	}

	var preassignToDepartmentId *uint = nil

	if mode == "reassign" {
		v, err := strconv.Atoi(*reassignToDepartmentId)

		if err != nil {
			return &out.CustomError{
				Code:    400,
				Message: "Invalid reassign_to_department_id value",
			}
		}

		preassignToDepartmentId = new(uint(v))
	}

	return useCase.storage.DeleteDepartment(ctx, uint(pId), mode, preassignToDepartmentId)
}
