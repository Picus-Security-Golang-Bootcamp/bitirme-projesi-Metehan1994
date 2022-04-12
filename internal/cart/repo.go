package cart

import (
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (c *CartRepository) Migration() {
	c.db.AutoMigrate(&models.Cart{})
	c.db.AutoMigrate(&models.CartItem{})
}

func (c *CartRepository) CreateBasket(userID uuid.UUID) {
	var cart models.Cart
	cart.UserID = userID
	result := c.db.Where("userID =?", cart.UserID).FirstOrCreate(&cart)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}

func (c *CartRepository) AddItem(CartID uuid.UUID) {
	var item models.CartItem
	item.CartID = CartID
	result := c.db.Where("cartID=?", item.CartID).FirstOrCreate(&item)
	if result.Error != nil {
		zap.L().Fatal(result.Error.Error())
	}
}
