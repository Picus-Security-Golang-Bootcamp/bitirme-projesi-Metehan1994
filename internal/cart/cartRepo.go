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

func (c *CartRepository) GetCartByUserID(UserID uuid.UUID) *models.Cart {
	var cart models.Cart
	cart.UserID = UserID
	c.db.Preload("Items.Product").Find(&cart)
	return &cart
}

func (c *CartRepository) GetOrCreateCart(userID uuid.UUID) (*models.Cart, string) {
	cart := c.GetCartByUserID(userID)
	var s string
	if cart.ID != uuid.Nil {
		s = "You have already a cart. New item will be added to it."
		fmt.Println(s)
	} else {
		cart.UserID = userID
		cart.ID = uuid.New()
		result := c.db.Where("user_id =?", cart.UserID).FirstOrCreate(&cart)
		if result.Error != nil {
			zap.L().Fatal(result.Error.Error())
		}
		fmt.Println("New cart is created for you.")
	}
	return cart, s
}

func (c *CartRepository) Update(cart *models.Cart) {
	result := c.db.Preload("Items.Product").Save(&cart)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}

func (c *CartRepository) List() models.Cart {
	var cart models.Cart
	c.db.Preload("Items.Product").Find(&cart)
	return cart
}

func (c *CartRepository) DeleteItemByID(cart *models.Cart, id int) error {
	err := c.cartItemRepo.DeleteById(uint(id))
	if err != nil {
		return err
	}
	c.Update(cart)
	return nil
}
