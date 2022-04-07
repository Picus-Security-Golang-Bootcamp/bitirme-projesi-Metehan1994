package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Price      int
	Quantity   int
	StockCode  int
	CategoryID int
	Category   Category `gorm:"OnDelete:SET NULL"`
}

func (Product) TableName() string {
	//default table name
	return "products"
}
