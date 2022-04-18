package order

import (
	"errors"
	"time"

	"github.com/Metehan1994/final-project/internal/cart"
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo     *OrderRepository
	productRepo   *product.ProductRepository
	orderItemRepo *OrderItemRepository
	cartRepo      *cart.CartRepository
}

func InitializeOrderService(orderRepo *OrderRepository, productRepo *product.ProductRepository, orderItemRepo *OrderItemRepository,
	cartRepo *cart.CartRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
	}
}

func (oserv *OrderService) GetCartByUserID(userID uuid.UUID) *models.Cart {
	cart := oserv.cartRepo.GetCartByUserID(userID)
	return cart
}

//CompleteOrder creates order and order item and update the product amount
func (oserv *OrderService) CompleteOrder(cart *models.Cart) (*models.Order, error) {
	tx := oserv.orderRepo.db.Begin()

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
		err := oserv.productRepo.UpdateProductQuantityAfterSale(&item)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
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
	}
	tx.Commit()
	return order, nil
}

func (oserv *OrderService) OrderGetWithItems(order *models.Order) (*models.Order, error) {
	orderWithItems, err := oserv.orderRepo.OrderGetWithItems(order)
	return orderWithItems, err
}

func (oserv *OrderService) DeleteCart(cart *models.Cart) error {
	oserv.cartRepo.DeleteCart(cart)
	return nil
}

func (oserv *OrderService) GetOrdersByUserID(userID uuid.UUID) ([]*models.Order, error) {
	orderList, err := oserv.orderRepo.GetOrdersByUserID(userID)
	return orderList, err
}

//CancelOrder manages the order cancelation by changing order status and updating product amount in the list
func (oserv *OrderService) CancelOrder(userID uuid.UUID, id string) (*models.Order, error) {
	var order *models.Order
	orders, err := oserv.GetOrdersByUserID(userID)
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
		err := oserv.orderRepo.UpdateOrder(order)
		if err != nil {
			return nil, err
		}
	}
	for _, item := range order.Items {
		err := oserv.productRepo.UpdateProductQuantityAfterCancel(&item)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}
