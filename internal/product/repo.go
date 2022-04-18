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

type IProductRepository interface {
	Create(a models.Product) (*models.Product, error)
	Update(models.Product) (*models.Product, error)
	DeleteById(id int) error
	GetByID(ID int) (*models.Product, error)
	GetBySku(sku string) *models.Product
	UpdateProductQuantityAfterSale(cartItem *models.CartItem) error
	UpdateProductQuantityAfterCancel(orderItem *models.OrderItem) error
}

func (p *ProductRepository) Migration() {
	p.db.AutoMigrate(&models.Product{})
}

func (p *ProductRepository) InsertSampleData(product *models.Product) {
	result := p.db.Unscoped().Where(models.Product{Name: product.Name}).FirstOrCreate(&product)
	if result.Error != nil {
		zap.L().Fatal("Cannot insert data into DB")
	}
}

//Create creates a new product
func (p *ProductRepository) Create(product models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("product", product))
	result := p.db.Where("sku = ?", product.Sku).FirstOrCreate(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	p.db.Preload("Category").Where("sku = ?", product.Sku).First(&product)
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
	zap.L().Debug("product.repo.getbyID", zap.Reflect("id", ID))
	var product models.Product
	result := p.db.Preload("Category").First(&product, ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &product, nil
}

//GetBySku provides the product info for a given sku
func (p *ProductRepository) GetBySku(sku string) *models.Product {
	var product *models.Product
	result := p.db.Where("sku=?", sku).First(&product)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return product
}

//UpdateProductQuantityAfterSale reduces the product amount which is found in the cart after completing order
func (p *ProductRepository) UpdateProductQuantityAfterSale(cartItem *models.CartItem) error {
	product := &cartItem.Product
	if product.Quantity < cartItem.Amount {
		return errors.New("not enough product. Please check product availability from the list once more")
	}
	result := p.db.Model(&product).Where("id = ? AND quantity >= ?", product.ID, cartItem.Amount).
		Update("quantity", gorm.Expr("quantity - ?", cartItem.Amount))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//UpdateProductQuantityAfterSale increases the product amount which is found in the order after canceling order
func (p *ProductRepository) UpdateProductQuantityAfterCancel(orderItem *models.OrderItem) error {
	product := &orderItem.Product
	result := p.db.Model(&product).Where("id = ?", product.ID).
		Update("quantity", gorm.Expr("quantity + ?", orderItem.Amount))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//ListProductsWithCategories provides a list of product with pagination
func (p *ProductRepository) ListProductsWithCategories(pageIndex, pageSize int) ([]*models.Product, int) {
	zap.L().Debug("product.repo.ListProducts")
	var allproducts []*models.Product
	p.db.Find(&allproducts)
	count := len(allproducts)
	var products []*models.Product
	p.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Preload("Category").Find(&products)
	return products, count
}

//SearchByName finds the products which includes the given word in name
func (p *ProductRepository) SearchByName(name string) ([]*models.Product, error) {
	var products []*models.Product
	result := p.db.Preload("Category").Where("name ILIKE ? ", "%"+name+"%").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

//SearchByName finds the products which includes the given word in sku
func (p *ProductRepository) SearchBySku(word string) ([]*models.Product, error) {
	var products []*models.Product
	result := p.db.Preload("Category").Where("sku ILIKE ? ", "%"+word+"%").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
