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
	c.db.Unscoped().Preload("Cart").Where("product_id=?", cartItem.ProductID).Create(&cartItem)
}

func (c *CartItemRepository) AddItem(cart *models.Cart, product *models.Product, quantity int) (*models.Cart, error) {
	var foundCart bool = false
	for _, item := range cart.Items {
		if item.ProductID == product.ID {
			foundCart = true
		}
	}
	if foundCart {
		return cart, errors.New("the product is already available in the basket. Please update its quantity")
	} else {
		cartItem := models.CartItem{
			ProductID: product.ID,
			Product:   *product,
			Price:     product.Price * quantity,
			Amount:    quantity,
			CartID:    cart.ID,
		}
		// cartItem.ProductID = product.ID
		// cartItem.Product = *product
		// cartItem.Price = product.Price * quantity
		// cartItem.Amount = quantity
		// cartItem.CartID = cart.ID
		c.CreateCartItem(&cartItem)
		cart.Items = append(cart.Items, cartItem)
		cart.TotalPrice += quantity * product.Price
		return cart, nil
	}
	// cartItem := c.GetItemByProductID(int(product.ID))
	// var cartItem2 models.CartItem
	// c.db.Unscoped().Where("product_id=?", product.ID).Find(&cartItem2)
	// if cartItem.ProductID == 0 && cartItem2.ProductID == 0 {
	// 	if product.Quantity < quantity {
	// 		return cart, errors.New("product amount is not enough to compensate your demand")
	// 	}
	// 	cartItem.ProductID = product.ID
	// 	cartItem.Product = *product
	// 	cartItem.Price = product.Price * quantity
	// 	cartItem.Amount = quantity
	// 	cartItem.CartID = cart.ID
	// 	c.CreateCartItem(cartItem)
	// 	cart.Items = append(cart.Items, *cartItem)
	// 	cart.TotalPrice += quantity * product.Price
	// 	return cart, nil
	// } else if cartItem.ProductID == 0 && cartItem2.ProductID != 0 {
	// 	if product.Quantity < quantity {
	// 		return cart, errors.New("product amount is not enough to compensate your demand")
	// 	}
	// 	cartItem.ProductID = product.ID
	// 	cartItem.Product = *product
	// 	cartItem.Price = product.Price * quantity
	// 	cartItem.Amount = quantity
	// 	cartItem.CartID = cart.ID
	// 	c.db.Unscoped().Preload("Cart").Create(cartItem)
	// 	cart.Items = append(cart.Items, *cartItem)
	// 	cart.TotalPrice += quantity * product.Price
	// 	return cart, nil
	// } else {
	// 	return cart, errors.New("the product is already available in the basket. Please update its quantity")
	// }
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
