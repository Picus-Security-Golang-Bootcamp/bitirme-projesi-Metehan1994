package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	Products    []Product
}

func (Category) TableName() string {
	//default table name
	return "categories"
}
