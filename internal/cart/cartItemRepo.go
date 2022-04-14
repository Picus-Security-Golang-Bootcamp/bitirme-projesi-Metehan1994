package cart

import (
	"errors"
	"fmt"

	"github.com/Metehan1994/final-project/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) *CartItemRepository {
	return &CartItemRepository{db: db}
}

func (c *CartItemRepository) Migration() {
	c.db.AutoMigrate(&models.CartItem{})
}

func (c *CartItemRepository) GetItemByProductID(productId int) *models.CartItem {
	var Item *models.CartItem
	c.db.Where("product_id=?", productId).Find(&Item)
	return Item
}

func (c *CartItemRepository) UpdateCartItem(cartItem *models.CartItem) {
	result := c.db.Preload("Product").Where("product_id=?", cartItem.ProductID).Save(&cartItem)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}

func (c *CartItemRepository) CreateCartItem(cartItem *models.CartItem) {
	c.db.Unscoped().Preload("Cart").Where("product_id=?", cartItem.ProductID).FirstOrCreate(&cartItem)
}

func (c *CartItemRepository) AddItem(cart *models.Cart, product *models.Product, quantity int) (*models.Cart, error) {
	cartItem := c.GetItemByProductID(int(product.ID))
	if cartItem.ProductID != 0 {
		return cart, errors.New("the product is already available in the basket. Please update its quantity")
	} else {
		if product.Quantity < quantity {
			return cart, errors.New("product amount is not enough to compensate your demand")
		}
		cartItem.ProductID = product.ID
		cartItem.Product = *product
		cartItem.Price = product.Price * quantity
		cartItem.Amount = quantity
		cartItem.CartID = cart.ID
		c.CreateCartItem(cartItem)
		cart.Items = append(cart.Items, *cartItem)
		cart.TotalPrice += quantity * product.Price
	}
	return cart, nil
}

//DeleteByID applies a soft delete to a cart item with given ID
func (c *CartItemRepository) DeleteById(id uint) error {
	zap.L().Debug("cartItem.repo.deleteById", zap.Reflect("id", id))
	var cartItem models.CartItem
	result := c.db.First(&cartItem, id)
	if result.Error != nil {
		return result.Error
	} else {
		fmt.Println("Valid ID, deleted:", id)
	}
	result = c.db.Delete(&models.CartItem{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *CartItemRepository) UpdateQuantityById(id int, quantity int) error {
	zap.L().Debug("cartItem.repo.updateQuantityById", zap.Reflect("id", id))
	var cartItem models.CartItem
	result := c.db.Preload("Product").First(&cartItem, id)
	if result.Error != nil {
		return result.Error
	}
	if cartItem.Product.Quantity < quantity {
		return errors.New("product quantity is not enough to compansate your demand")
	}
	cartItem.Amount = quantity
	cartItem.Price = cartItem.Product.Price * quantity
	c.UpdateCartItem(&cartItem)
	return nil
}
