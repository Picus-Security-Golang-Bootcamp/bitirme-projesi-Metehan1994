package order

import (
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) Migration() {
	o.db.AutoMigrate(&models.Order{})
}

//OrderGetWithItems list the order items in an order with products
func (o *OrderRepository) OrderGetWithItems(order *models.Order) (*models.Order, error) {
	zap.L().Debug("Order.repo.OrderGetWithItems", zap.Reflect("order", order))
	result := o.db.Preload("Items.Product").First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

//OrderGetWithItems list the number of orders created by a user
func (o *OrderRepository) GetOrdersByUserID(userID uuid.UUID) ([]*models.Order, error) {
	zap.L().Debug("Order.repo.OrderGetByUserID", zap.Reflect("userID", userID))
	var orders []*models.Order
	result := o.db.Preload("Items.Product").Where("user_id=?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (o *OrderRepository) UpdateOrder(order *models.Order) error {
	result := o.db.Where("id=?", order.ID).Save(&order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
