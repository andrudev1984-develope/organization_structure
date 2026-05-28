package out

type Department struct {
	Id        uint
	ParentId  *uint
	Name      string
	CreatedAt string
}

type Employee struct {
	Id        uint
	FullName  string
	Position  string
	CreatedAt string
	HiredAt   *string
}

type DepartmentSingleInfo struct {
	Department Department
	Employees  []Employee
	Children   []Department
}

type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}
