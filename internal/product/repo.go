package product

import (
	"errors"
	"fmt"

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

//Create creates a new product
func (p *ProductRepository) Create(product models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("product", product))
	result := p.db.Where("sku = ?", product.Sku).FirstOrCreate(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

//Update updates the info for the product attributes
func (p *ProductRepository) Update(product models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.update", zap.Reflect("product", product))
	result := p.db.Preload("Category").Save(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

//DeleteByID applies a soft delete to a product with given ID
func (p *ProductRepository) DeleteById(id int) error {
	zap.L().Debug("product.repo.deleteById", zap.Reflect("id", id))
	var product models.Product
	result := p.db.First(&product, id)
	if result.Error != nil {
		return result.Error
	} else {
		fmt.Println("Valid ID, deleted:", id)
	}
	result = p.db.Delete(&models.Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//GetByID provides the product info for a given ID
func (p *ProductRepository) GetByID(ID int) (*models.Product, error) {
	var product models.Product
	result := p.db.Preload("Category").First(&product, ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &product, nil
}

func (p *ProductRepository) UpdateProductQuantityAfterSale(cartItem *models.CartItem) error {
	product := &cartItem.Product
	fmt.Println(product)
	if product.Quantity < cartItem.Amount {
		return errors.New("not enough product. Please check product availability from the list once more")
	}
	result := p.db.Model(&product).Where("id = ? AND quantity >= ?", product.ID, cartItem.Amount).
		Update("quantity", gorm.Expr("quantity - ?", cartItem.Amount))
	if result.Error != nil {
		return result.Error
	}
	fmt.Println(product.Quantity)
	return nil
}

func (p *ProductRepository) UpdateProductQuantityAfterCancel(orderItem *models.OrderItem) error {
	product := &orderItem.Product
	fmt.Println(product)
	result := p.db.Model(&product).Where("id = ?", product.ID).
		Update("quantity", gorm.Expr("quantity + ?", orderItem.Amount))
	if result.Error != nil {
		return result.Error
	}
	fmt.Println(product.Quantity)
	return nil
}
