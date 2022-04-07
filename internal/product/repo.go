package product

import (
	"github.com/Metehan1994/final-project/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Migration() {
	r.db.AutoMigrate(&models.Product{})
}
