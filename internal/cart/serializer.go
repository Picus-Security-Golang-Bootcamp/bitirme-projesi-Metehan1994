package cart

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
)

func CartToResponse(p *models.Cart) *api.Cart {
	items := make([]*api.CartItem, 0)
	for _, item := range p.Items {
		items = append(items, CartItemToResponse(&item))
	}
	status := p.GetCartStatusAsString()
	return &api.Cart{
		UserID:     p.UserID.String(),
		ID:         p.ID.String(),
		Status:     status,
		TotalPrice: int64(p.TotalPrice),
		Items:      items,
	}
}

func CartItemToResponse(p *models.CartItem) *api.CartItem {
	return &api.CartItem{
		ProductName: p.Product.Name,
		ID:          int64(p.ID),
		Amount:      int64(p.Amount),
		Price:       int64(p.Price),
	}
}
