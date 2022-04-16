package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/Metehan1994/final-project/internal/models"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db            *gorm.DB
	orderItemRepo *OrderItemRepository
	productRepo   *product.ProductRepository
}

func NewOrderRepository(db *gorm.DB, orderItemRepo *OrderItemRepository, productRepo *product.ProductRepository) *OrderRepository {
	return &OrderRepository{db: db, orderItemRepo: orderItemRepo, productRepo: productRepo}
}

func (o *OrderRepository) Migration() {
	o.db.AutoMigrate(&models.Order{})
}

//CompleteOrder creates order and order item and update the product amount
func (o *OrderRepository) CompleteOrder(cart *models.Cart) (*models.Order, error) {
	zap.L().Debug("Order.repo.CompleteOrder", zap.Reflect("cart", cart))
	tx := o.db.Begin()

	order := &models.Order{
		ID:          uuid.New(),
		CartID:      cart.ID,
		UserID:      cart.UserID,
		OrderStatus: 1,
		TotalPrice:  cart.TotalPrice,
	}
	result := tx.Create(order)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	for _, item := range cart.Items {
		orderItem := &models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Price:     item.Price,
			Amount:    item.Amount,
		}
		result := tx.Create(orderItem)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
		err := o.productRepo.UpdateProductQuantityAfterSale(&item)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return order, nil
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

//CancelOrder manages the order cancelation by changing order status and updating product amount in the list
func (o *OrderRepository) CancelOrder(userID uuid.UUID, id string) (*models.Order, error) {
	zap.L().Debug("Order.repo.CancelOrder")
	var order *models.Order
	orders, err := o.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}
	var orderFound bool = false
	for _, ord := range orders {
		if ord.ID.String() == id {
			orderFound = true
			order = ord
		}
	}
	if !orderFound {
		return nil, errors.New("no available order with this ID")
	} else if order.GetOrderStatusAsString() == "canceled" {
		return nil, errors.New("it has been already canceled")
	} else if time.Now().After(order.CreatedAt.AddDate(0, 0, 14)) {
		return nil, errors.New("you cannot cancel the order after 14 days")
	} else {
		order.OrderStatus = 2
		result := o.db.Where("id=?", order.ID).Save(&order)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	for _, item := range order.Items {
		fmt.Println(item.Product) //Productı buluyor.
		err := o.productRepo.UpdateProductQuantityAfterCancel(&item)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}
