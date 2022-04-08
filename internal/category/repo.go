package category

import (
	"github.com/Metehan1994/final-project/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (c *CategoryRepository) Migration() {
	c.db.AutoMigrate(&models.Category{})
}

func (c *CategoryRepository) InsertSampleData(category *models.Category) models.Category {
	result := c.db.Unscoped().Where(models.Category{Name: category.Name}).FirstOrCreate(&category)
	if result.Error != nil {
		zap.L().Fatal("Cannot insert data into DB") //Check Error
	}

	return *category
}
