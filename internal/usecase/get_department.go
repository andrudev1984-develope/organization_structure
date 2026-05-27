package usecase

import (
	"context"
	"organization_structure/internal/model/out"
	"strconv"
	"time"
)

func (useCase *UseCase) GetDepartment(ctx context.Context, id string, depth string, includeEmployees string) (out.DepartmentSingleInfo, error) {
	if depth == "" {
		depth = "1"
	}

	if includeEmployees == "" {
		includeEmployees = "true"
	}

	if _, err := strconv.Atoi(id); err != nil {
		return out.DepartmentSingleInfo{}, &out.CustomError{
			Code:    400,
			Message: "Need valid department id",
		}
	}

	if _, err := strconv.Atoi(depth); err != nil {
		return out.DepartmentSingleInfo{}, &out.CustomError{
			Code:    400,
			Message: "Need valid depth",
		}
	}

	pId, _ := strconv.Atoi(id)
	pdepth, err := strconv.Atoi(depth)

	if pdepth > 5 {
		return out.DepartmentSingleInfo{}, &out.CustomError{
			Code:    400,
			Message: "Need valid depth : 1-5",
		}
	}

	pIncludeEmployees, err := strconv.ParseBool(includeEmployees)

	if err != nil {
		return out.DepartmentSingleInfo{}, &out.CustomError{
			Code:    400,
			Message: "Need valid includeEmployees : true or false",
		}
	}

	dep, err := useCase.storage.GetDepartment(ctx, uint(pId), pdepth, pIncludeEmployees)

	if err != nil {
		return out.DepartmentSingleInfo{}, err
	}

	defer func() {
		dep.Employees = nil
		dep.Children = nil
	}()

	var info out.DepartmentSingleInfo

	info.Department.Id = dep.Id
	info.Department.Name = dep.Name
	info.Department.CreatedAt = dep.CreatedAt.Format("02-01-2006")

	info.Employees = make([]out.Employee, 0, len(dep.Employees))
	info.Children = make([]out.Department, 0, len(dep.Children))

	for _, e := range dep.Employees {
		eo := out.Employee{
			Id:        e.Id,
			FullName:  e.FullName,
			Position:  e.Position,
			CreatedAt: e.CreatedAt.Format("02-01-2006"),
			HiredAt:   formatDate((*time.Time)(e.HiredAt)),
		}
		info.Employees = append(info.Employees, eo)
	}

	for _, d := range dep.Children {
		do := out.Department{
			Id:        d.Id,
			CreatedAt: d.CreatedAt.Format("02-01-2006"),
			Name:      d.Name,
		}
		info.Children = append(info.Children, do)
	}

	return info, nil
}
