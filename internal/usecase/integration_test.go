package usecase

import (
	"context"
	"log/slog"
	"organization_structure/internal/adapters/database"
	"organization_structure/internal/adapters/orm"
	"organization_structure/internal/config"
	"organization_structure/internal/model/in"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	g := orm.SetupDatabase(config.Config{
		Db: struct {
			Host     string
			Port     int
			User     string
			Password string
			Name     string
			SslMode  string
		}{
			Host:     "postgres_storage_test",
			Port:     5432,
			User:     "postgres",
			Password: "12345678",
			Name:     "organization_structure",
			SslMode:  "disable",
		},
		Server: struct{ Port int }{
			Port: 0,
		},
		Log: struct{ Level string }{
			Level: "",
		},
	})

	adapter := database.NewPgStorage(g)
	var usecase = NewUseCase(adapter)

	defer func() {
		d, _ := g.DB()

		d.Close()
	}()

	t.Run("TestCreateDepartment", func(t *testing.T) {
		root, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		slog.Error("Parent created", slog.Any("id", root.Id), slog.Any("CreatedAt", root.CreatedAt))

		_, e = usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: new(uint(0)),
		})

		assert.Error(t, e)

		slog.Error("Error", slog.Any("error", e))

		child, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: new(root.Id),
		})

		assert.Nil(t, e)

		slog.Info("Child created", slog.Any("id", child.Id), slog.Any("CreatedAt", child.CreatedAt),
			slog.Any("Name", child.Name), slog.Any("ParentId", *child.ParentId))

		_, e = usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     child.Name,
			ParentId: child.ParentId,
		})

		assert.Error(t, e)

		slog.Error("Error", slog.Any("error", e))

		slog.Info("Success create department")
	})

	t.Run("TestCreateEmployee", func(t *testing.T) {
		root, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		_, e = usecase.CreateEmployee(context.Background(), strconv.Itoa(0), in.CreateEmployee{
			FullName: "Name",
			Position: "Position",
		})

		assert.Error(t, e)

		slog.Error("Error", slog.Any("error", e))

		emp, e := usecase.CreateEmployee(context.Background(), strconv.Itoa(int(root.Id)), in.CreateEmployee{
			FullName: "Name",
			Position: "Position",
		})

		assert.Nil(t, e)

		slog.Info("Employee created", slog.Any("id", emp.Id), slog.Any("CreatedAt", emp.CreatedAt),
			slog.Any("Full name", emp.FullName), slog.Any("Position", emp.Position))

		assert.Equal(t, "Name", emp.FullName)
		assert.Equal(t, "Position", emp.Position)

		slog.Info("Success create employee")
	})

	t.Run("TestGetDepartment", func(t *testing.T) {
		root, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		child, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: new(root.Id),
		})

		assert.Nil(t, e)

		emp, e := usecase.CreateEmployee(context.Background(), strconv.Itoa(int(root.Id)), in.CreateEmployee{
			FullName: "Name",
			Position: "Position",
		})

		assert.Nil(t, e)

		gD, e := usecase.GetDepartment(context.Background(), strconv.Itoa(int(root.Id)), strconv.Itoa(1), "true")

		assert.Nil(t, e)

		assert.Equal(t, root.Id, gD.Department.Id)

		assert.Equal(t, 1, len(gD.Children))
		assert.Equal(t, 1, len(gD.Employees))

		assert.Equal(t, child.Id, gD.Children[0].Id)
		assert.Equal(t, emp.Id, gD.Employees[0].Id)
		assert.Equal(t, emp.FullName, gD.Employees[0].FullName)
		assert.Equal(t, emp.Position, gD.Employees[0].Position)

		slog.Info("Success get department")
	})

	t.Run("TestUpdateDepartment", func(t *testing.T) {
		root, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		root1, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		d, e := usecase.UpdateDepartment(context.Background(), strconv.Itoa(int(root1.Id)), in.UpdateDepartment{
			Name:     new(root1.Name + "Updated"),
			ParentId: new(root.Id),
		})

		assert.Nil(t, e)

		assert.Equal(t, root1.Name+"Updated", d.Department.Name)
		assert.Equal(t, root.Id, *d.Department.ParentId)

		_, e = usecase.UpdateDepartment(context.Background(), strconv.Itoa(int(root.Id)), in.UpdateDepartment{
			Name:     nil,
			ParentId: new(root1.Id),
		})

		assert.Error(t, e)

		slog.Error("Error", slog.Any("error", e))

		slog.Info("Success update department")
	})

	t.Run("TestDeleteDepartment", func(t *testing.T) {
		root, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		child1, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: new(root.Id),
		})

		assert.Nil(t, e)

		child2, e := usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: new(child1.Id),
		})

		assert.Nil(t, e)

		e = usecase.DeleteDepartment(context.Background(), strconv.Itoa(int(root.Id)), "cascade", nil)

		assert.Nil(t, e)

		_, e = usecase.GetDepartment(context.Background(), strconv.Itoa(int(root.Id)), strconv.Itoa(1), "false")

		assert.Error(t, e)

		_, e = usecase.GetDepartment(context.Background(), strconv.Itoa(int(child1.Id)), strconv.Itoa(1), "false")

		assert.Error(t, e)

		_, e = usecase.GetDepartment(context.Background(), strconv.Itoa(int(child2.Id)), strconv.Itoa(1), "false")

		assert.Error(t, e)

		root, e = usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: nil,
		})

		assert.Nil(t, e)

		child1, e = usecase.CreateDepartment(context.Background(), in.CreateDepartment{
			Name:     "Name" + time.Now().String(),
			ParentId: new(root.Id),
		})

		assert.Nil(t, e)

		_, e = usecase.CreateEmployee(context.Background(), strconv.Itoa(int(root.Id)), in.CreateEmployee{
			FullName: "Name",
			Position: "Position",
		})

		e = usecase.DeleteDepartment(context.Background(), strconv.Itoa(int(root.Id)), "reassign", new(strconv.Itoa(int(child1.Id))))

		assert.Nil(t, e)

		_, e = usecase.GetDepartment(context.Background(), strconv.Itoa(int(root.Id)), strconv.Itoa(1), "false")

		assert.Error(t, e)

		child1M, e := usecase.GetDepartment(context.Background(), strconv.Itoa(int(child1.Id)), strconv.Itoa(1), "true")

		assert.Equal(t, 1, len(child1M.Employees))

		slog.Info("Success delete department")
	})
}
