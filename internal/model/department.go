package model

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	Id        uint `gorm:"primaryKey"`
	ParentID  *uint
	Name      string
	Parent    *Department
	Children  []Department `gorm:"foreignkey:ParentID"`
	Employees []Employee
	DeletedAt gorm.DeletedAt
	CreatedAt time.Time
}
