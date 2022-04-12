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

func (c *CategoryRepository) createCategory(category *models.Category) (*models.Category, error) {
	zap.L().Debug("author.repo.create", zap.Reflect("categoryBody", category))
	if err := c.db.Create(category).Error; err != nil {
		zap.L().Error("author.repo.Create failed to create category", zap.Error(err))
		return nil, err
	}

	return category, nil
}

func (c *CategoryRepository) ListCategoriesWithProducts() ([]models.Category, error) {
	zap.L().Debug("author.repo.ListCategories")
	var categories []models.Category
	result := c.db.Preload("Products").Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
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
