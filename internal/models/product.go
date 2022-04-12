package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	Price       int
	Quantity    int
	StockCode   string
	CategoryID  uint
	Category    Category `gorm:"OnDelete:SET NULL"`
}

func (Product) TableName() string {
	//default table name
	return "products"
}
