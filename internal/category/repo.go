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

//ListCategoriesWithProducts brings the list of categories with their products
func (c *CategoryRepository) ListCategoriesWithProducts(pageIndex, pageSize int) ([]*models.Category, int) {
	zap.L().Debug("category.repo.ListCategories")
	var allcategories []*models.Category
	c.db.Find(&allcategories)
	count := len(allcategories)
	var categories []*models.Category
	c.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Preload("Products").Find(&categories)
	return categories, count
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
