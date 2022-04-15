package order

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
)

func OrderToResponse(p *models.Order) *api.Order {
	items := make([]*api.OrderItem, 0)
	for _, item := range p.Items {
		items = append(items, OrderItemToResponse(&item))
	}
	status := p.GetOrderStatusAsString()
	return &api.Order{
		UserID:     p.UserID.String(),
		ID:         p.ID.String(),
		Status:     status,
		TotalPrice: int64(p.TotalPrice),
		Items:      items,
	}
}

func OrderListToResponse(p []*models.Order) []*api.Order {
	orders := make([]*api.Order, 0)
	for _, order := range p {
		orders = append(orders, OrderToResponse(order))
	}
	return orders
}

func OrderItemToResponse(p *models.OrderItem) *api.OrderItem {
	return &api.OrderItem{
		ProductName: p.Product.Name,
		ProductID:   int64(p.ProductID),
		ID:          int64(p.ID),
		Price:       int64(p.Price),
		Amount:      int64(p.Amount),
	}
}
