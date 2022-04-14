package order

import (
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db            *gorm.DB
	orderItemRepo *OrderItemRepository
}

func NewOrderRepository(db *gorm.DB, orderItemRepo *OrderItemRepository) *OrderRepository {
	return &OrderRepository{db: db, orderItemRepo: orderItemRepo}
}

func (o *OrderRepository) Migration() {
	o.db.AutoMigrate(&models.Order{})
}

func (o *OrderRepository) CompleteOrder(cart *models.Cart) (*models.Order, error) {
	tx := o.db.Begin()

	order := &models.Order{
		ID:          uuid.New(),
		CartID:      cart.ID,
		UserID:      cart.UserID,
		OrderStatus: 1,
		TotalPrice:  cart.TotalPrice,
	}
	result := o.db.Create(order)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	for _, item := range cart.Items {
		orderItem := &models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Price:     item.Price,
		}
		result := o.db.Create(orderItem)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
	}
	tx.Commit()
	return order, nil
}

func (o *OrderRepository) OrderGetWithItems(order *models.Order) (*models.Order, error) {
	result := o.db.Preload("Items.Product").First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (o *OrderRepository) GetOrdersByUserID(userID uuid.UUID) ([]*models.Order, error) {
	var orders []*models.Order
	result := o.db.Preload("Items.Product").Where("user_id=?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
