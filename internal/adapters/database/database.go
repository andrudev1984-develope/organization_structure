package database

import (
	"context"
	"errors"
	"organization_structure/internal/model"
	"organization_structure/internal/model/out"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PgStorage struct {
	adapter *gorm.DB
}

const (
	GetChildDepartments = `WITH RECURSIVE cte_name (id, parent_id, name) AS (
    select id, parent_id, name, created_at, 0 as depth from departments where id = $1
    UNION ALL
    select d.id, d.parent_id, d.name, d.created_at, c.depth + 1 from departments d inner join cte_name c on d.parent_id = c.id and depth < $2)
	SELECT id, parent_id, name, created_at FROM cte_name where depth > 0;`

	GetChildCycle = `WITH RECURSIVE cte_name (id, parent_id, name) AS (
    select id, parent_id, name, 0 as depth from departments where id = $1
    UNION ALL
    select d.id, d.parent_id, d.name, c.depth + 1 from departments d inner join cte_name c on d.parent_id = c.id)
    cycle id set is_cycle to true default false using path_nodes
	select exists (SELECT * FROM cte_name where is_cycle is true);`
)

func (s PgStorage) CreateDepartment(ctx context.Context, name string, parentId *uint) (*model.Department, error) {
	d := model.Department{
		ParentID: parentId,
		Name:     name,
	}

	if parentId != nil {
		_, err := gorm.G[model.Department](s.adapter).Where("id=?", *parentId).First(ctx)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, &out.CustomError{Code: 404, Message: "Parent department is not found"}
			}

			return nil, err
		}
	}

	err := gorm.G[model.Department](s.adapter).Create(ctx, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (s PgStorage) GetDepartment(ctx context.Context, departmentId uint, depth int, includeEmployees bool) (*model.Department, error) {
	pre := gorm.G[model.Department](s.adapter).Where("id=?", departmentId)

	if includeEmployees {
		pre = pre.Preload("Employees", nil)
	}

	dep, err := pre.First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &out.CustomError{Code: 404, Message: "Department is not found"}
		}

		return nil, err
	}

	deps, err := gorm.G[model.Department](s.adapter).Raw(GetChildDepartments, departmentId, depth).Find(ctx)

	defer func() {
		deps = nil
	}()

	if err != nil {
		return nil, &out.CustomError{Code: 500, Message: err.Error()}
	}

	for _, d := range deps {
		dep.Children = append(dep.Children, d)
	}

	return &dep, nil
}

func (s PgStorage) UpdateDepartment(ctx context.Context, departmentId uint, name *string, parentId *uint) error {
	dep, err := gorm.G[model.Department](s.adapter).Where("id=?", departmentId).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &out.CustomError{Code: 404, Message: "Department is not found"}
		}

		return err
	}

	if parentId != nil {
		_, err = gorm.G[model.Department](s.adapter).Where("id=?", *parentId).First(ctx)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &out.CustomError{Code: 404, Message: "Parent department is not found"}
			}
			return err
		}
	}

	var data = make(map[string]interface{})

	if name != nil {
		data["name"] = *name
	}

	if parentId == nil {
		data["parent_id"] = nil
	} else {
		data["parent_id"] = *parentId
	}

	err = s.adapter.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&dep).Updates(data)

		if err.Error != nil {
			return err.Error
		}

		res, errc := gorm.G[bool](tx).Raw(GetChildCycle, departmentId).First(ctx)

		if errc != nil {
			return errc
		}

		if res {
			return &out.CustomError{Code: 409, Message: "Find cycle in department hierarchy"}
		}

		return nil
	})

	return err
}

func (s PgStorage) DeleteDepartment(ctx context.Context, departmentId uint, mode string, reassignToDepartmentId *uint) error {
	pre := gorm.G[model.Department](s.adapter).Where("id=?", departmentId)

	if mode == "reassign" {
		pre = pre.Preload("Employees", nil)
	}

	dep, err := pre.First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &out.CustomError{Code: 404, Message: "Department is not found"}
		}

		return err
	}

	if mode == "cascade" {
		// delete all child departments / employees
		s.adapter.Unscoped().Delete(&dep)
	} else {
		err = s.adapter.Transaction(func(tx *gorm.DB) error {
			// sort delete department
			_, err := gorm.G[model.Department](tx).Where("id = ?", departmentId).Delete(ctx)

			if err != nil {
				return err
			}

			var ids = make([]uint, 0, len(dep.Employees))

			defer func() {
				ids = nil
			}()

			for _, v := range dep.Employees {
				ids = append(ids, v.Id)
			}

			// move its employees to reassignToDepartmentId
			_, err = gorm.G[model.Employee](tx).Where("id in ?", ids).Updates(ctx, model.Employee{DepartmentID: *reassignToDepartmentId})

			return err
		})

	}

	return err
}

func (s PgStorage) CreateEmployee(ctx context.Context, fullName string, departmentId uint, position string, hiredAt *time.Time) (*model.Employee, error) {
	e := model.Employee{
		DepartmentID: departmentId,
		FullName:     fullName,
		Position:     position,
		HiredAt:      (*datatypes.Date)(hiredAt),
	}

	_, err := gorm.G[model.Department](s.adapter).Where("id=?", departmentId).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &out.CustomError{Code: 404, Message: "Department is not found"}
		}

		return nil, err
	}

	err = gorm.G[model.Employee](s.adapter).Create(ctx, &e)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func NewPgStorage(adapter *gorm.DB) *PgStorage {
	return &PgStorage{adapter: adapter}
}
