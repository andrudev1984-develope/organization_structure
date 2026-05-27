package in

type CreateDepartment struct {
	Name     string
	ParentId *uint
}

type UpdateDepartment struct {
	Name     *string
	ParentId *uint
}

type CreateEmployee struct {
	FullName string
	Position string
	HiredAt  *string
}
