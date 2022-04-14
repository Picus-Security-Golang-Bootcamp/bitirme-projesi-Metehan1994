package order

import (
	"github.com/Metehan1994/final-project/internal/models"
	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (o *OrderItemRepository) Migration() {
	o.db.AutoMigrate(&models.OrderItem{})
}
