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

//GetItemByProductID finds the cart item by product ID
func (c *CartItemRepository) GetItemByProductID(productId int) *models.CartItem {
	zap.L().Debug("cartItem.repo.getItemById", zap.Reflect("productid", productId))
	var Item *models.CartItem
	c.db.Where("product_id=?", productId).Find(&Item)
	return Item
}

//UpdateCartItem updates the cart item based on the changes
func (c *CartItemRepository) UpdateCartItem(cartItem *models.CartItem) {
	zap.L().Debug("cartItem.repo.UpdateCartItem", zap.Reflect("cartItem", cartItem))
	result := c.db.Preload("Product").Where("product_id=?", cartItem.ProductID).Save(&cartItem)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}

//CreateCartItem creates a new cart item
func (c *CartItemRepository) CreateCartItem(cartItem *models.CartItem) {
	zap.L().Debug("cartItem.repo.CreateCartItem", zap.Reflect("cartItem", cartItem))
	result := c.db.Unscoped().Preload("Cart").Where("product_id=?", cartItem.ProductID).Create(&cartItem)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}

//DeleteByID applies a soft delete to a cart item with given ID
func (c *CartItemRepository) DeleteById(cart *models.Cart, id uint) (*models.Cart, error) {
	zap.L().Debug("cartItem.repo.deleteById", zap.Reflect("id", id))
	var cartItem models.CartItem
	result := c.db.First(&cartItem, id)
	if result.Error != nil {
		return nil, result.Error
	} else {
		fmt.Println("Valid ID, deleted:", id)
		cart.TotalPrice -= cartItem.Price
	}
	result = c.db.Delete(&models.CartItem{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
}

//UpdateQuantity updates the quantity of cart item with given ID
func (c *CartItemRepository) UpdateQuantityById(cart *models.Cart, id int, quantity int) (*models.Cart, error) {
	zap.L().Debug("cartItem.repo.updateQuantityById", zap.Reflect("id", id))
	var cartItem models.CartItem
	result := c.db.Preload("Product").First(&cartItem, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if cartItem.Product.Quantity < quantity {
		return nil, errors.New("product quantity is not enough to compansate your demand")
	}
	cartItem.Amount = quantity
	cart.TotalPrice -= cartItem.Price
	cartItem.Price = cartItem.Product.Price * quantity
	cart.TotalPrice += cartItem.Price
	c.UpdateCartItem(&cartItem)
	return cart, nil
}
