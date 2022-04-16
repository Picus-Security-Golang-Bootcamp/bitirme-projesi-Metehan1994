package cart

import (
	"fmt"

	"github.com/Metehan1994/final-project/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartRepository struct {
	db           *gorm.DB
	cartItemRepo *CartItemRepository
}

func NewCartRepository(db *gorm.DB, cartItemRepo *CartItemRepository) *CartRepository {
	return &CartRepository{db: db, cartItemRepo: cartItemRepo}
}

func (c *CartRepository) Migration() {
	c.db.AutoMigrate(&models.Cart{})
}

//GetCartByUserID finds the cart beloging to a user by id
func (c *CartRepository) GetCartByUserID(UserID uuid.UUID) *models.Cart {
	zap.L().Debug("cart.repo.GetCartByUserID", zap.Reflect("UserID", UserID))
	var cart models.Cart
	//cart.UserID = UserID
	c.db.Preload("Items.Product").Where("user_id=?", UserID).Find(&cart)
	fmt.Println(cart)
	return &cart
}

//GerOrCreateCart checks the cart is available or not and creates cart if user does not have a cart
func (c *CartRepository) CreateCart(cart *models.Cart) *models.Cart {
	zap.L().Debug("cart.repo.GetOrCreateCart", zap.Reflect("cart", cart))
	cart.ID = uuid.New()
	result := c.db.Where("user_id =?", cart.UserID).FirstOrCreate(&cart)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
	fmt.Println("New cart is created for you.")
	return cart
}

//Update updates the changes done on the cart
func (c *CartRepository) Update(cart *models.Cart) {
	zap.L().Debug("cart.repo.Update", zap.Reflect("cart", cart))
	result := c.db.Preload("Items.Product").Where("id=?", cart.ID).Save(&cart)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}

//List shows the cart of a user
func (c *CartRepository) List() models.Cart {
	zap.L().Debug("cart.repo.List")
	var cart models.Cart
	c.db.Preload("Items.Product").Where("user_id=?").Find(&cart)
	return cart
}

//DeleteCart removes the cart of user
func (c *CartRepository) DeleteCart(cart *models.Cart) error {
	zap.L().Debug("cart.repo.deleteById", zap.Reflect("cart", cart))
	result := c.db.First(&cart, cart.ID)
	if result.Error != nil {
		return result.Error
	} else {
		fmt.Println("Valid ID, deleted:", cart.ID)
	}
	for _, item := range cart.Items {
		_, err := c.cartItemRepo.DeleteById(cart, item.ID)
		if err != nil {
			return err
		}
	}
	result = c.db.Delete(&models.Cart{}, cart.ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
