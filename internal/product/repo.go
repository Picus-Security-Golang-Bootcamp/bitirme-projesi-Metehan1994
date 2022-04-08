package product

import (
	"github.com/Metehan1994/final-project/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) Migration() {
	p.db.AutoMigrate(&models.Product{})
}

func (p *ProductRepository) InsertSampleData(product *models.Product) {
	result := p.db.Unscoped().Where(models.Product{Name: product.Name}).FirstOrCreate(&product)
	if result.Error != nil {
		zap.L().Fatal("Cannot insert data into DB") //Check error
	}
}
