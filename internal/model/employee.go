package model

import (
	"time"

	"gorm.io/datatypes"
)

type Employee struct {
	Id           uint `gorm:"primaryKey"`
	DepartmentID uint
	FullName     string
	Position     string
	Department   Department
	HiredAt      *datatypes.Date
	CreatedAt    time.Time
}
